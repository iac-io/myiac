package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dfernandezm/myiac/internal/util"
	"google.golang.org/api/iam/v1"
	"io"
	"os"
	"strings"
)

// IamClient wrapper interface around GCP Service Account Client that allows creation and listing of Service Account keys.
// Instead of extracting the exact methods from GCP *iam.Service, they are wrapped on these 2 operations (CreateKey, ListKeys).
// This is due to the complicated call chains inside *iam.Service (i.e. 'Projects.ServiceAccounts.Keys.Create(...)')
type IamClient interface {
	CreateKey(request *iam.CreateServiceAccountKeyRequest, resource string) (*iam.ServiceAccountKey, error)
	ListKeys(resource string) (*iam.ListServiceAccountKeysResponse, error)
}

// gcpIamClient private implementation of the above interface that actually wraps the operations
type gcpIamClient struct {
	iamService *iam.Service
	ctx context.Context
}

// NewGcpIamClient Creates a new IamClient with external provided context and *iam.Service
func NewGcpIamClient(ctx context.Context, iamService *iam.Service) *gcpIamClient {
	return &gcpIamClient{iamService: iamService, ctx: ctx}
}

// NewGcpIamClient Creates a new IamClient with default context.
// Authentication against GCP must have already been performed when invoking this operation
func NewDefaultIamClient() *gcpIamClient {
	ctx := context.Background()
	return NewGcpIamClient(ctx, getIamService(ctx))
}

func (gic *gcpIamClient) CreateKey(request *iam.CreateServiceAccountKeyRequest, resource string) (*iam.ServiceAccountKey, error) {
	return gic.iamService.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
}

func (gic *gcpIamClient) ListKeys(resource string) (*iam.ListServiceAccountKeysResponse, error) {
	return gic.iamService.Projects.ServiceAccounts.Keys.List(resource).Do()
}

func getIamService(ctx context.Context) *iam.Service {
	service, err := iam.NewService(ctx)
	if err != nil {
		panic(err)
	}
	return service
}

// ObjectStorageClient is the interface for the client.Storage GCP SDK operations so they can be mocked in tests
// or injected
type ObjectStorageClient interface {
	Bucket(bucketName string) *storage.BucketHandle
}

// ObjectStorageCache allows management of basic object storage operations in GCP (write, read).
// Abstracts the operations needed from the GCP SDK
type ObjectStorageCache interface {
	Write(ctx context.Context, bucketName string, key string, content interface{}) error
	Read(ctx context.Context, bucketName string, objectKey string) (interface{},error)
}

const keysCacheBucketName = "moneycol-keys"

// gcpObjectStorageCache Implements 'ObjectStorageCache' and uses the interface
// 'ObjectStorageClient' to perform the base operations
type gcpObjectStorageCache struct {
	client ObjectStorageClient
	ctx context.Context
}

// NewGcpObjectStorageCache creates a GCP-based Object Storage cache, optionally
// receiving a Context. It injects a ObjectStorageClient providing the base operations
func NewGcpObjectStorageCache(ctx context.Context, client ObjectStorageClient) *gcpObjectStorageCache {
	return &gcpObjectStorageCache{client:client, ctx:ctx}
}

// NewDefaultObjectStorageCache creates a GCP-based Object Storage cache using a
// inner context
func NewDefaultObjectStorageCache() *gcpObjectStorageCache {
	ctx := context.Background()
	return NewGcpObjectStorageCache(ctx, getObjectStorageClient(ctx))
}

// Write writes the given 'objectContent' into a bucket called 'bucketName' with
// key 'key'. If a Context is passed, it's used instead of the default one to
// regenerate the GCP SDK client
func (gosc *gcpObjectStorageCache) Write(ctx context.Context, bucketName string, key string, objectContent interface{}) error {
	// replace context with provided
	if ctx != nil {
		gosc.client = getObjectStorageClient(ctx)
	} else {
		ctx = context.Background()
	}

	bkt := gosc.client.Bucket(bucketName)
	obj := bkt.Object(key)

	// Write something to obj.
	// w implements io.Writer.
	w := obj.NewWriter(ctx)

	// Write some text to obj. This will either create the object or overwrite whatever is there already.
	if _, err := fmt.Fprintf(w, "%v", objectContent); err != nil {
		return err
	}

	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

// Read performs a 'read' from bucket 'bucketName' on key 'key'. A passed Context could be used
// to regenerate the client
func (gosc *gcpObjectStorageCache) Read(ctx context.Context, bucketName string, key string) (interface{},error)  {
	if ctx != nil {
		gosc.client = getObjectStorageClient(ctx)
	} else {
		ctx = context.Background()
	}

	bkt := gosc.client.Bucket(bucketName)
	obj := bkt.Object(key)

	//TODO: should list the bucket first, see:
	// https://godoc.org/cloud.google.com/go/storage
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("getting object error %v", err)
	}

	defer r.Close()

	buf := new(strings.Builder)

	if _, copyErr := io.Copy(buf, r); copyErr != nil {
		return "", fmt.Errorf("error copying to stdout %v", copyErr)
	}

	keyString := buf.String()
	fmt.Printf("Found value %s\n", keyString)
	return keyString, nil
}

// ServiceAccountClient allows management of service account keys (create, list) for given service account emails
type ServiceAccountClient interface {
	KeyForServiceAccount(saEmail string, recreateKey bool) (string, error)
	CreateKey(serviceAccountEmail string)
	ListKeys(serviceAccountEmail string) ([]string, error)
}

type serviceAccountClient struct {
	gcpIamClient       IamClient
	objectStorageCache ObjectStorageCache
}

// NewServiceAccountClient creates a new GCP client for service account key management
func NewServiceAccountClient(iamClient IamClient, objectStorageCache ObjectStorageCache) *serviceAccountClient {
	return &serviceAccountClient{gcpIamClient: iamClient, objectStorageCache: objectStorageCache}
}

// NewDefaultServiceAccountClient creates a new GCP client for service account key management based on defaults
// Authentication against GCP must have already been performed before invoking this operation
func NewDefaultServiceAccountClient() *serviceAccountClient {
	return NewServiceAccountClient(NewDefaultIamClient(), NewDefaultObjectStorageCache())
}

// CreateKey creates a service account key for the given service account email
func (sac *serviceAccountClient) CreateKey(serviceAccountEmail string) (string, error) {
	request := &iam.CreateServiceAccountKeyRequest{}
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail)
	fmt.Printf("Creating key\n")
	key, err :=  sac.gcpIamClient.CreateKey(request, resource)

	if err != nil {
		return "", fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Created key: %v", key.Name)
	fmt.Printf("Created the key %v\n", key)
	return key.PrivateKeyData, nil
}

func (sac *serviceAccountClient) ListKeys(serviceAccountEmail string) ([]string, error) {
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail)
	response, err := sac.gcpIamClient.ListKeys(resource)

	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.List: %v", err)
	}

	res := make([]string,0)
	for _, key := range response.Keys {
		res = append(res, key.PrivateKeyData)
	}

	return res, nil
}

// KeyForServiceAccount given a service account email address, it creates or recreates a JSON key for it which is
// returned as a string
func (sac *serviceAccountClient) KeyForServiceAccount(saEmail string, recreateKey bool) (string, error) {
	keys, _ := sac.ListKeys(saEmail)
	var jsonKey string

	if len(keys) > 0 && !recreateKey {
		fmt.Printf("Using existing key for SA: %s\n", saEmail)
		keyString, err := sac.objectStorageCache.Read(nil, keysCacheBucketName, saEmail+".json")
		if err != nil {
			return "", fmt.Errorf("err finding key %v", err)
		}
		jsonKey = keyString.(string)
	} else {
		newKey, err := sac.CreateKey(saEmail)
		if err != nil {
			fmt.Printf("Error creating key %v", err)
			return "", err
		}

		jsonKeyString := util.Base64Decode(newKey)

		fmt.Println("----- BEGIN Service account JSON key -------")
		fmt.Println(jsonKeyString)
		fmt.Println("----- END Service account JSON key ---------")

		writeErr := sac.objectStorageCache.Write(nil, keysCacheBucketName, saEmail+".json", jsonKeyString)
		if writeErr != nil {
			return "", fmt.Errorf("error writing key to storage %v", writeErr)
		}

		jsonKey = jsonKeyString
	}

	return jsonKey, nil
}

func getObjectStorageClient(ctx context.Context) ObjectStorageClient {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return client
}

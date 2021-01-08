package gcp

//TODO: move this into gcp package as subpackage?
import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iam/v1"
)

const KeysCacheBucketName = "moneycol-keys"

// IamClient wrapper interface around GCP Service Account Client that allows creation and listing of Service Account keys.
// Instead of extracting the exact methods from GCP *iam.Service, they are wrapped on these 2 operations (CreateKey, ListKeys).
// This is due to the complicated call chains inside *iam.Service (i.e. 'Projects.ServiceAccounts.Keys.Create(...)')
type IamGcpClient interface {
	CreateKey(request *iam.CreateServiceAccountKeyRequest, resource string) (*iam.ServiceAccountKey, error)
	ListKeys(resource string) (*iam.ListServiceAccountKeysResponse, error)
}

// gcpIamClient private implementation of the above interface that actually wraps the operations
type gcpIamClient struct {
	iamService *iam.Service
	ctx        context.Context
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
type ObjectStorageGcpClient interface {
	Bucket(bucketName string) *storage.BucketHandle
}

// ObjectStorageCache allows management of basic object storage operations in GCP (write, read).
// Abstracts the operations needed from the GCP SDK
type ObjectStorageCache interface {
	Write(ctx context.Context, bucketName string, key string, content interface{}) error
	Read(ctx context.Context, bucketName string, objectKey string) (interface{}, error)
}

// gcpObjectStorageCache Implements 'ObjectStorageCache' and uses the interface
// 'ObjectStorageClient' to perform the base operations
type gcpObjectStorageCache struct {
	client ObjectStorageGcpClient
	ctx    context.Context
}

// NewGcpObjectStorageCache creates a GCP-based Object Storage cache, optionally
// receiving a Context. It injects a ObjectStorageClient providing the base operations
func NewGcpObjectStorageCache(ctx context.Context, client ObjectStorageGcpClient) *gcpObjectStorageCache {
	return &gcpObjectStorageCache{client: client, ctx: ctx}
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
func (gosc *gcpObjectStorageCache) Read(ctx context.Context, bucketName string, key string) (interface{}, error) {
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

func getObjectStorageClient(ctx context.Context) ObjectStorageGcpClient {
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return client
}

package gcp

import (
	"fmt"
	"os"
	"strings"

	"github.com/iac-io/myiac/internal/util"
	"google.golang.org/api/iam/v1"
)

// ServiceAccountClient allows management of service account keys (create, list) for given service account emails
type ServiceAccountClient interface {
	KeyForServiceAccount(saEmail string, recreateKey bool) (string, error)
	KeyFileForServiceAccount(saEmail string, recreateKey bool, filePath string) error
	CreateKey(serviceAccountEmail string)
	ListKeys(serviceAccountEmail string) ([]string, error)
}

type serviceAccountClient struct {
	gcpIamClient       IamGcpClient
	objectStorageCache ObjectStorageCache
}

// NewServiceAccountClient creates a new GCP client for service account key management
func NewServiceAccountClient(iamClient IamGcpClient, objectStorageCache ObjectStorageCache) *serviceAccountClient {
	return &serviceAccountClient{gcpIamClient: iamClient, objectStorageCache: objectStorageCache}
}

// NewDefaultServiceAccountClient creates a new GCP client for service account key management based on defaults
// Authentication against GCP must have already been performed before invoking this operation
func NewDefaultServiceAccountClient() *serviceAccountClient {
	return NewServiceAccountClient(NewDefaultIamClient(), NewDefaultObjectStorageCache())
}

// CreateKey creates a service account key for the given service account email
func (sac *serviceAccountClient) CreateKey(serviceAccountEmail string) (string, string, error) {
	request := &iam.CreateServiceAccountKeyRequest{}
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail)
	fmt.Printf("Creating key...\n")
	key, err := sac.gcpIamClient.CreateKey(request, resource)

	if err != nil {
		return "", "", fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Created key: %v\n", key.Name)
	return key.PrivateKeyData, key.Name, nil
}

func (sac *serviceAccountClient) ListKeys(serviceAccountEmail string) ([]string, error) {
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccountEmail)
	response, err := sac.gcpIamClient.ListKeys(resource)

	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.List: %v", err)
	}

	var keyIds []string
	for _, key := range response.Keys {
		keyName := key.Name
		keyId := extractKeyId(keyName)
		keyIds = append(keyIds, keyId)
	}

	return keyIds, nil
}

// KeyForServiceAccount given a service account email address, it creates or recreates a JSON key for it which is
// returned back as a string
func (sac *serviceAccountClient) KeyForServiceAccount(saEmail string, recreateKey bool) (string, error) {
	keyIds, _ := sac.ListKeys(saEmail)
	keyJsonFile := saEmail + ".json"

	if len(keyIds) > 0 && !recreateKey {
		// Even though keys may exist already for the service account, they cannot be returned if they aren't cached
		// at the time they have been created
		fmt.Printf("Attempt to find existing key for SA: %s\n", saEmail)

		for _, keyId := range keyIds {
			fmt.Printf("Checking if keyId %s has been cached\n", keyId)
			cachedKeyJsonFile := addKeyIdToJsonFile(keyJsonFile, keyId)
			keyString, err := sac.objectStorageCache.Read(nil, KeysCacheBucketName, cachedKeyJsonFile)
			if err == nil {
				fmt.Printf("Returning key for SA: %s from cache\n", saEmail)
				return keyString.(string), nil
			}
		}

		fmt.Printf("Could not find cached key for SA: %s -- Creating new one\n", saEmail)
	}

	return sac.createAndCacheKey(saEmail, keyJsonFile)
}

// KeyFileForServiceAccount same as 'KeyForServiceAccount' but returns the key written to 'filePath' destination
func (sac *serviceAccountClient) KeyFileForServiceAccount(saEmail string, recreateKey bool, filePath string) error {
	jsonKey, err := sac.KeyForServiceAccount(saEmail, recreateKey)

	if err != nil {
		fmt.Printf("Error generating Key for SA %s %v", saEmail, err)
		return err
	}

	writeErr := util.WriteStringToFile(jsonKey, filePath)

	if writeErr != nil {
		fmt.Printf("Error writing Key to file %v\n", writeErr)
		return writeErr
	}

	return nil
}

func (sac *serviceAccountClient) createAndCacheKey(saEmail string, keyJsonFile string) (string, error) {
	keyData, keyName, err := sac.CreateKey(saEmail)
	if err != nil {
		fmt.Printf("Error creating key %v\n", err)
		return "", err
	}

	jsonKeyString := util.Base64Decode(keyData)

	fmt.Println("----- BEGIN Service account JSON key -------")
	fmt.Println(jsonKeyString)
	fmt.Println("----- END Service account JSON key ---------")
	keyId := extractKeyId(keyName)
	keyJsonFile = addKeyIdToJsonFile(keyJsonFile, keyId)
	fmt.Printf("writing with cache key: %s\n", keyJsonFile)
	writeErr := sac.objectStorageCache.Write(nil, KeysCacheBucketName, keyJsonFile, jsonKeyString)
	if writeErr != nil {
		return "", fmt.Errorf("error caching key %v", writeErr)
	}
	return jsonKeyString, nil
}

func extractKeyId(keyName string) string {
	idx := strings.LastIndex(keyName, "/")
	keyId := keyName[idx+1:]
	fmt.Printf("KeyName: %s -> KeyId %s\n", keyName, keyId)
	return keyId
}

func addKeyIdToJsonFile(keyJsonFile string, keyId string) string {
	return strings.Replace(keyJsonFile, ".json", "-"+keyId+".json", -1)
}

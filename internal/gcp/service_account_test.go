package gcp

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

type keyGenerator struct {
	keyName string
	keyData string
	keys    []*iam.ServiceAccountKey
}

func newKeyGenerator(keyName string, keyData string) *keyGenerator {
	return &keyGenerator{keyName: keyName, keyData: keyData}
}

func (kg *keyGenerator) keyGenerator() *iam.ServiceAccountKey {
	key := &iam.ServiceAccountKey{
		KeyAlgorithm:    "",
		KeyOrigin:       "",
		KeyType:         "",
		Name:            kg.keyName,
		PrivateKeyData:  kg.keyData,
		PrivateKeyType:  "",
		PublicKeyData:   "",
		ValidAfterTime:  "",
		ValidBeforeTime: "",
		ServerResponse:  googleapi.ServerResponse{},
		ForceSendFields: nil,
		NullFields:      nil,
	}

	kg.keys = append(kg.keys, key)
	return key
}

type fakeIamClient struct {
	*keyGenerator
}

func newFakeIamClient(keyGenerator *keyGenerator) *fakeIamClient {
	return &fakeIamClient{keyGenerator: keyGenerator}
}

func (fic *fakeIamClient) setKeyGenerator(keyGenerator *keyGenerator) {
	fic.keyGenerator = keyGenerator
}

func (fic fakeIamClient) CreateKey(request *iam.CreateServiceAccountKeyRequest, resource string) (*iam.ServiceAccountKey, error) {
	key := fic.keyGenerator.keyGenerator()
	return key, nil
}

func (fic fakeIamClient) ListKeys(resource string) (*iam.ListServiceAccountKeysResponse, error) {
	resp := &iam.ListServiceAccountKeysResponse{
		Keys:            fic.keyGenerator.keys,
		ServerResponse:  googleapi.ServerResponse{},
		ForceSendFields: nil,
		NullFields:      nil,
	}
	return resp, nil
}

// fakeObjectStorageCache see mock creation using testify / mockery
// https://blog.lamida.org/mocking-in-golang-using-testify/ and
// https://tutorialedge.net/golang/improving-your-tests-with-testify-go/
type fakeObjectStorageCache struct {
	mock.Mock
}

func (m *fakeObjectStorageCache) Write(ctx context.Context, bucketName string, key string, content interface{}) error {
	args := m.Called(ctx, bucketName, key, content)
	return args.Error(0)
}

func (m fakeObjectStorageCache) Read(ctx context.Context, bucketName string, objectKey string) (interface{}, error) {
	args := m.Called(ctx, bucketName, objectKey)
	return args.Get(0).(interface{}), args.Error(1)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreateNewKey(t *testing.T) {
	kg := newKeyGenerator("testKey", "testData")
	mockIamClient := newFakeIamClient(kg)

	objStorageCacheMock := new(fakeObjectStorageCache)
	objStorageCacheMock.On("Write", nil, "test", "key", "val").Return(nil)
	objStorageCacheMock.On("Read", nil, "testBucket", "key").Return("", nil)

	saClient := NewServiceAccountClient(mockIamClient, objStorageCacheMock)

	testSaEmail := "testAccount@gcloudserviceaccount.com"
	keyData, keyName, _ := saClient.CreateKey(testSaEmail)

	// cache isn't used when creating keys directly
	objStorageCacheMock.AssertNumberOfCalls(t, "Write", 0)
	objStorageCacheMock.AssertNumberOfCalls(t, "Read", 0)

	assert.Equal(t, "testKey", keyName)
	assert.Equal(t, "testData", keyData)

}

package gcp

import (
	"context"
	"fmt"
	"log"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/iac-io/myiac/internal/encryption"
	"github.com/iac-io/myiac/internal/util"
	"google.golang.org/api/iterator"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type keyRing struct {
	name string
}

type key struct {
	name string
}

type kmsEncrypter struct {
	projectId          string
	locationId         string
	defaultKeyRingName string
	defaultKeyName     string
	keyRingHolder      *keyRing
	keyHolder          *key
	kmsClient          *kms.KeyManagementClient
}

// Creates a new encrypter using GCP KMS service
//
// See: https://cloud.google.com/kms/docs/creating-keys
//
func NewKmsEncrypter(projectId string, locationId string, defaultKeyRingName string,
	defaultKeyName string) encryption.Encrypter {
	kenc := new(kmsEncrypter)
	kenc.projectId = projectId
	kenc.locationId = locationId
	kenc.defaultKeyName = defaultKeyName
	kenc.defaultKeyRingName = defaultKeyRingName
	return kenc
}

func (kenc *kmsEncrypter) Encrypt(plainText string) (string, error) {
	ctx := context.Background()
	client, _ := kenc.createKmsClient(ctx)

	// Convert the message into bytes. Cryptographic plaintext and
	// ciphertext are always byte arrays.
	plaintext := []byte(plainText)

	// Get or create a keyRing
	if kenc.keyRingHolder == nil {
		keyRing, err := kenc.getOrCreateKeyRing()
		if err != nil {
			//TODO: should use stacktraces
			log.Fatalf("Error creating keyRing %v", err)
		}
		kenc.keyRingHolder = keyRing
	}

	// Get or create a key
	if kenc.keyHolder == nil {
		key, err := kenc.getOrCreateCryptoKey(kenc.keyRingHolder.name, kenc.defaultKeyName)
		if err != nil {
			log.Fatalf("Error creating key %v", err)
		}
		kenc.keyHolder = key
	}

	// Build the encrypt request.
	req := &kmspb.EncryptRequest{
		Name:      kenc.keyHolder.name,
		Plaintext: plaintext,
	}

	result, err := client.Encrypt(ctx, req)
	if err != nil {
		log.Fatalf("failed to encrypt: %v", err)
	}

	cipherText := string(result.Ciphertext)
	fmt.Printf("Encrypted ciphertext: %s\n", cipherText)

	//TODO: convert to base64? before returning?
	return cipherText, nil
}

func (kenc *kmsEncrypter) Decrypt(cipherText string) (string, error) {
	ctx := context.Background()
	client, _ := kenc.createKmsClient(ctx)

	// Convert the message into bytes. Cryptographic plaintext and
	// ciphertext are always byte arrays.
	ciphertext := []byte(cipherText)

	// Obtain the key
	kenc.getEncryptionKey()

	// Build the encrypt request.
	req := &kmspb.DecryptRequest{
		Name:       kenc.keyHolder.name,
		Ciphertext: ciphertext,
	}

	result, err := client.Decrypt(ctx, req)
	if err != nil {
		log.Fatalf("failed to decrypt: %v", err)
	}

	plainText := string(result.Plaintext)
	fmt.Printf("Encrypted plainText: %s\n", plainText)

	//TODO: convert back from base64? before returning?
	return plainText, nil
}

func (kenc *kmsEncrypter) getEncryptionKey() {
	// Get or create a keyRing
	if kenc.keyRingHolder == nil {
		keyRing, err := kenc.getOrCreateKeyRing()
		if err != nil {
			//TODO: should use stacktraces
			log.Fatalf("Error creating keyRing %v", err)
		}
		kenc.keyRingHolder = keyRing
	}

	// Get or create a key
	if kenc.keyHolder == nil {
		key, err := kenc.getOrCreateCryptoKey(kenc.keyRingHolder.name, kenc.defaultKeyName)
		if err != nil {
			log.Fatalf("Error creating key %v", err)
		}
		kenc.keyHolder = key
	}
}

// createKeyRing creates a new ring to store keys on KMS.
// parent := "projects/PROJECT_ID/locations/global"
// id := "my-key-ring"
func (kenc *kmsEncrypter) createKeyRing(id string, parent string) (*keyRing, error) {

	// Create the client.
	ctx := context.Background()
	client, _ := kenc.createKmsClient(ctx)

	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: id,
	}

	result, err := client.CreateKeyRing(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to create key ring: %v", err)
	}

	fmt.Printf("Created key ring: %s\n", result.Name)
	fullKeyRingName := fmt.Sprintf("%s/keyRings/%s", parent, result.Name)
	keyRing := &keyRing{name: fullKeyRingName}
	return keyRing, nil
}

func (kenc *kmsEncrypter) listKeyRings(parent string) ([]string, error) {
	// Create the client.
	ctx := context.Background()
	client, _ := kenc.createKmsClient(ctx)

	listReq := &kmspb.ListKeyRingsRequest{
		Parent: parent,
	}

	result := client.ListKeyRings(ctx, listReq)
	var keyRings []string
	for {
		resp, err := result.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list key rings: %v", err)
		}

		fmt.Printf("key ring: %s\n", resp.Name)
		keyRings = append(keyRings, resp.Name)
	}

	return keyRings, nil
}

func (kenc *kmsEncrypter) getOrCreateKeyRing() (*keyRing, error) {
	parent := fmt.Sprintf("projects/%s/locations/%s",
		kenc.projectId,
		kenc.locationId)

	id := kenc.defaultKeyRingName
	keyRings, _ := kenc.listKeyRings(parent)

	fmt.Printf("Check if keyRing with id %s exists\n", id)
	fullKeyRingName := fmt.Sprintf("%s/keyRings/%s", parent, id)

	if util.ArrayContains(keyRings, fullKeyRingName) {
		fmt.Printf("Keyring exists with id %s, using it\n", id)
		existingKeyRing := &keyRing{name: fullKeyRingName}
		return existingKeyRing, nil
	} else {
		fmt.Printf("Keyring with id %s does not exist, creating now\n", id)
		newKeyRing, err := kenc.createKeyRing(id, parent)

		if err != nil {
			return nil, fmt.Errorf("failed to create keyring: %v\n", err)
		}

		return newKeyRing, nil
	}
}

// Creates or retrieves an existing symmetric encryption key
func (kenc *kmsEncrypter) getOrCreateCryptoKey(parentKeyRing string, keyName string) (*key, error) {

	fullKeyName := fmt.Sprintf("%s/cryptoKeys/%s",
		parentKeyRing,
		keyName)
	keys, _ := kenc.listCryptoKeys(parentKeyRing)

	fmt.Printf("Found keys %v\n", keys)
	fmt.Printf("Check if key with id %s exists\n", keyName)

	if util.ArrayContains(keys, fullKeyName) {
		fmt.Printf("Keyring exists with id %s, using it\n", keyName)
		existingKeyRing := &key{name: fullKeyName}
		return existingKeyRing, nil
	} else {
		createdKeyName, err := kenc.createKeySymmetricEncryptDecrypt(parentKeyRing, keyName)

		if err != nil {
			return nil, fmt.Errorf("failed to create key: %v\n", err)
		}

		fmt.Printf("Created key with name %s", createdKeyName)
		newKey := &key{name: createdKeyName}
		return newKey, nil
	}
}

// createKeySymmetricEncryptDecrypt creates a new symmetric encrypt/decrypt key
// on Cloud KMS.
//
// parentKeyRing := "projects/my-project/locations/us-east1/keyRings/my-key-ring"
// id := "my-symmetric-encryption-key"
func (kenc *kmsEncrypter) createKeySymmetricEncryptDecrypt(parentKeyRing string, id string) (string, error) {

	// Create the client.
	ctx := context.Background()
	client, _ := kenc.createKmsClient(ctx)

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parentKeyRing,
		CryptoKeyId: id,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
		},
	}

	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create key: %v", err)
	}
	fmt.Printf("Created key: %s\n", result.Name)
	return result.Name, nil
}

func (kenc *kmsEncrypter) listCryptoKeys(parent string) ([]string, error) {
	// Create the client.
	ctx := context.Background()
	client, err := kenc.createKmsClient(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %v", err)
	}

	listReq := &kmspb.ListCryptoKeysRequest{
		Parent: parent,
	}

	result := client.ListCryptoKeys(ctx, listReq)
	var keys []string
	for {
		resp, err := result.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list keys: %v", err)
		}

		fmt.Printf("crypto key: %s\n", resp.Name)
		keys = append(keys, resp.Name)
	}

	return keys, nil
}

func (kenc *kmsEncrypter) createKmsClient(ctx context.Context) (*kms.KeyManagementClient, error) {

	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %v", err)
	}
	return client, nil
}

//TODO: add a key rotation period

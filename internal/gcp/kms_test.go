package gcp

import (
	"log"
	"os"
	"testing"

	"github.com/iac-io/myiac/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateGCPKMSService(t *testing.T) {
	t.Skip("GCP KMS setup is required for this test")
	gcpClient := NewKmsEncrypter("moneycol", "moneycol-keyring",
		"moneycol-keyring", "moneycol-infra-key")
	assert.NotNil(t, gcpClient)
}

func TestEncrypts(t *testing.T) {
	t.Skip("GCP KMS setup is required for this test")
	homeDir, _ := os.UserHomeDir()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", homeDir+"/moneycol_account.json")
	gcpClient := NewKmsEncrypter("moneycol", "global", "moneycol-keyring",
		"moneycol-infra-key")
	cipherText, err := gcpClient.Encrypt("a very sensitive value")
	log.Printf("Test ciphertext: %s\n", cipherText)
	_ = util.WriteStringToFile(cipherText, "/tmp/encrypted-value-2.txt")
	assert.Equal(t, nil, err)
}

func TestDecrypts(t *testing.T) {
	t.Skip("GCP KMS setup is required for this test")
	homeDir, _ := os.UserHomeDir()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", homeDir+"/moneycol_account.json")
	gcpClient := NewKmsEncrypter("moneycol", "global", "moneycol-keyring",
		"moneycol-infra-key")
	cipherToDecrypt, _ := util.ReadFileToString("/tmp/encrypted-value-2.txt")
	plainText, err := gcpClient.Decrypt(cipherToDecrypt)
	log.Printf("Test plainText: %s\n", plainText)
	assert.Equal(t, nil, err)
}

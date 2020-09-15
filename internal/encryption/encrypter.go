package encryption

import "github.com/dfernandezm/myiac/internal/util"

type Encrypter interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

type encryptionService struct {
	encrypter Encrypter
}
// NewEncryptionService create a new encrypter
func NewEncrypter(encrypter Encrypter) *encryptionService {
	enc := new(encryptionService)
	enc.encrypter = encrypter
	return enc
}

// EncryptFileContents encrypts the text contained in 'filename' and returns cipherText
// into a file with the same name ended with '.enc'
func (enc encryptionService) EncryptFileContents(filename string) string {
	plainText, _ := util.ReadFileToString(filename)
	cipherText, _ := enc.encrypter.Encrypt(plainText)
	cipherTextFilename := filename + ".enc"
	_ = util.WriteStringToFile(cipherText, cipherTextFilename)
	return cipherTextFilename
}

// DecryptFileContents decrypts the ciphertext contained in 'filename' and returns the
// plaintext result into another file with same name ended with '.dec'
func (enc encryptionService) DecryptFileContents(filename string) string {
	cipherText, _ := util.ReadFileToString(filename)
	plainText, _ := enc.encrypter.Decrypt(cipherText)
	plainTextFilename := filename + ".dec"
	_ = util.WriteStringToFile(plainText, plainTextFilename)
	return plainTextFilename
}

package ssl

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/secret"
	"github.com/dfernandezm/myiac/app/util"
	"log"
)

type CertStore interface {
	Register(certificate *Certificate)
}

type SecretCertStore struct {
	secretManager secret.SecretManager
}

func NewSecretCertStore(secretManager secret.SecretManager) *SecretCertStore {
	return &SecretCertStore{secretManager:secretManager}
}

func (scs *SecretCertStore) Register(certificate *Certificate) {
	tlsSecret := secret.NewTlsSecret(certificate.Domain, certificate.certPath, certificate.privateKeyPath)
	scs.secretManager.CreateTlsSecret(tlsSecret)
}

type Certificate struct {
	Domain string
	privateKey string
	privateKeyPath string
	cert string
	certPath string
}

// NewCertificate Create a new certificate from paths to PEM and KEY files
func NewCertificate(domainName string, certPath string, privateKeyPath string) *Certificate {
	certValue, privateKeyValue := certValuesFromPaths(certPath, privateKeyPath)
	return &Certificate{
		Domain:domainName,
		cert: certValue,
		privateKey: privateKeyValue,
		certPath: certPath,
		privateKeyPath: privateKeyPath}
}

// Read certificate values as strings from given file paths
func certValuesFromPaths(certPath string, privateKeyPath string) (string, string) {
	certValue, err := util.ReadFileToString(certPath)

	if err != nil {
		fmt.Printf("error reading cert %v", err)
		log.Fatal(err)
	}

	privateKeyValue, err2 := util.ReadFileToString(privateKeyPath)

	if err2 != nil {
		fmt.Printf("error reading key %v", err)
		log.Fatal(err2)
	}

	return certValue, privateKeyValue
}
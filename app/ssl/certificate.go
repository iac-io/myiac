package ssl

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/util"
	"log"
)

type Certificate struct {
	Domain string
	caBundle string
	privateKey string
	privateKeyPath string
	cert string
	certPath string
}

// NewCertificateFromLocation Create a new certificate from paths to PEM and KEY files
func NewCertificateFromLocation(domainName string, certPath string, privateKeyPath string) *Certificate {
	certValue, privateKeyValue := certValuesFromPaths(certPath, privateKeyPath)
	return &Certificate{
		Domain:domainName,
		cert: certValue,
		privateKey: privateKeyValue,
		certPath: certPath,
		privateKeyPath: privateKeyPath}
}

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
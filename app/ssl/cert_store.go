package ssl

import (
	"github.com/dfernandezm/myiac/app/cluster"
	"os"
)

type CertStore interface {
	Save(certificate *Certificate)
}

// create tls.cert / tls.key
type KubernetesSecretStore struct {
	namespace string
	secretName string
}

func NewKubernetesSecretStore(namespace string) *KubernetesSecretStore {
	kst := new(KubernetesSecretStore)
	kst.namespace = namespace
	return kst
}

// Store stores the Certificate in a Kubernetes secret
//
// The secretName is the domainName
func (kst *KubernetesSecretStore) Save(certificate *Certificate) {
	kst.secretName = certificate.Domain
	tlsKeyPath := "/tmp/tls.key"
	tlsCertPath := "/tmp/tls.crt"
	_ = os.Rename(certificate.privateKeyPath, tlsKeyPath)
	_ = os.Rename(certificate.certPath, tlsCertPath)
	cluster.CreateTlsSecret(kst.secretName, kst.namespace, tlsCertPath, tlsKeyPath)
}
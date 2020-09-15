package secret

import (
	"github.com/dfernandezm/myiac/internal/cluster"
	"os"
)

type SecretManager interface {
	CreateTlsSecret(secret TlsSecret)
}

type TlsSecret struct {
	name string
	tlsCertPath string
	tlsKeyPath string
}

func NewTlsSecret(name string, tlsCertPath string, tlsKeyPath string) TlsSecret {
	return TlsSecret{
		name: name,
		tlsCertPath: tlsCertPath,
		tlsKeyPath:  tlsKeyPath,
	}
}

type kubernetesSecretManager struct {
	namespace        string
	kubernetesRunner cluster.KubernetesRunner
}

func NewKubernetesSecretManager(namespace string, kuberneetesRunner cluster.KubernetesRunner) *kubernetesSecretManager {
	return &kubernetesSecretManager{
		namespace: namespace,
		kubernetesRunner:kuberneetesRunner,
	}
}

func (ksm kubernetesSecretManager) CreateTlsSecret(secret TlsSecret) {
	tlsKeyPathTmp := "/tmp/tls.key"
	tlsCertPathTmp := "/tmp/tls.crt"
	_ = os.Rename(secret.tlsKeyPath, tlsKeyPathTmp)
	_ = os.Rename(secret.tlsCertPath, tlsCertPathTmp)
	ksm.kubernetesRunner.CreateTlsSecret(secret.name, ksm.namespace, tlsKeyPathTmp, tlsCertPathTmp)
}

func (ksm kubernetesSecretManager) FindTlsSecret(secretName string) {
	ksm.kubernetesRunner.FindSecret(secretName, ksm.namespace)
}
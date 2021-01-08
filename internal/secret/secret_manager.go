package secret

import (
	"os"

	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/commandline"
)

type SecretManager interface {
	CreateTlsSecret(secret TlsSecret)
}

type TlsSecret struct {
	name        string
	tlsCertPath string
	tlsKeyPath  string
}

func NewTlsSecret(name string, tlsCertPath string, tlsKeyPath string) TlsSecret {
	return TlsSecret{
		name:        name,
		tlsCertPath: tlsCertPath,
		tlsKeyPath:  tlsKeyPath,
	}
}

type kubernetesSecretManager struct {
	namespace        string
	kubernetesRunner cluster.KubernetesRunner
}

func NewKubernetesSecretManager(namespace string, kubernetesRunner cluster.KubernetesRunner) *kubernetesSecretManager {
	return &kubernetesSecretManager{
		namespace:        namespace,
		kubernetesRunner: kubernetesRunner,
	}
}

func CreateKubernetesSecretManager(namespace string) *kubernetesSecretManager {
	return NewKubernetesSecretManager(namespace, cluster.NewKubernetesRunner(commandline.NewEmpty()))
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

func (ksm kubernetesSecretManager) CreateFileSecret(secretName string, filePath string) {
	ksm.kubernetesRunner.CreateFileSecret(secretName, ksm.namespace, filePath)
}

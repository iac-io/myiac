package secret

import (
	"os"

	"github.com/iac-io/myiac/internal/commandline"
)

const tlsKeyPathTmp = "/tmp/tls.key"
const tlsCertPathTmp = "/tmp/tls.crt"

type SecretManager interface {
	CreateTlsSecret(secret TlsSecret)
	CreateFileSecret(secretName string, filePath string)
	CreateLiteralSecret(secretName string, literalsMap map[string]string)
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
	kubernetesRunner KubernetesSecretRunner
}

func NewKubernetesSecretManager(namespace string, kubernetesRunner KubernetesSecretRunner) SecretManager {
	return &kubernetesSecretManager{
		namespace:        namespace,
		kubernetesRunner: kubernetesRunner,
	}
}

func CreateKubernetesSecretManager(namespace string) SecretManager {
	return NewKubernetesSecretManager(namespace, NewKubernetesRunner(commandline.NewEmpty()))
}

func (ksm kubernetesSecretManager) CreateTlsSecret(secret TlsSecret) {
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

func (ksm kubernetesSecretManager) CreateLiteralSecret(secretName string, literalsMap map[string]string) {
	ksm.kubernetesRunner.CreateLiteralSecret(secretName, ksm.namespace, literalsMap)
}

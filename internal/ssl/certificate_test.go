package ssl

import (
	"github.com/dfernandezm/myiac/internal/cluster"
	"github.com/dfernandezm/myiac/internal/commandline"
	"github.com/dfernandezm/myiac/internal/secret"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeKubernetesRunner struct {
}

func (fk fakeKubernetesRunner) SetupWithoutOutput(cmd string, args []string) {

}

func (fk fakeKubernetesRunner) Run() commandline.CommandOutput {
 return commandline.CommandOutput{Output: "Success"}
}



func TestCreateTlsCertificate(t *testing.T) {
	// setup
	cmdLine := commandline.NewEmpty()
	kubernetesRunner := cluster.NewKubernetesRunner(cmdLine)
	secretManager := secret.NewKubernetesSecretManager("default", kubernetesRunner)

	// given
	domain := "test-domain"
	certPath := "/Users/david/Documents/cloudflare-cert/cert.pem"
	keyPath := "/Users/david/Documents/cloudflare-cert/cert.key"

	//_ = util.WriteStringToFile("testCert", certPath)
	//_ = util.WriteStringToFile("testKey", keyPath)

	// when
	certificate := NewCertificate(domain, certPath, keyPath)
	certStore := NewSecretCertStore(secretManager)
	certStore.Register(certificate)

	// then
	createdSecretName := kubernetesRunner.FindSecret(domain, "default")
	assert.Contains(t, createdSecretName, domain)
}


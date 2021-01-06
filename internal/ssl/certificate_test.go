package ssl

import (
	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/secret"
	"github.com/iac-io/myiac/internal/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type fakeKubernetesRunner struct {
	cmd string
	args []string
	CmdLines []string
}

func (fk *fakeKubernetesRunner) SetupWithoutOutput(cmd string, args []string) {
	fk.cmd = cmd
	fk.args = args
}

func (fk *fakeKubernetesRunner) Run() commandline.CommandOutput {
	currentCmdLine := fk.cmd + " " + strings.Join(fk.args, " ")
	fk.CmdLines = append(fk.CmdLines, currentCmdLine)
	return commandline.CommandOutput{Output: "test-domain"}
}

func (fk fakeKubernetesRunner) RunVoid() {
}

func (fk fakeKubernetesRunner) Output() string {
	return "test-domain"
}

func (fk fakeKubernetesRunner) Setup(cmd string, args []string) {
}

func (fk fakeKubernetesRunner) IgnoreError(ignoreError bool) {
}

func (fk fakeKubernetesRunner) SetupCmdLine(cmdLine string) {
}


func TestCreateTlsCertificate(t *testing.T) {
	// setup
	cmdLine := new(fakeKubernetesRunner)
	kubernetesRunner := cluster.NewKubernetesRunner(cmdLine)
	secretManager := secret.NewKubernetesSecretManager("default", kubernetesRunner)

	// given
	domain := "test-domain"
	certPath := "/tmp/cert.pem"
	keyPath := "/tmp/cert.key"

	_ = util.WriteStringToFile("testCert", certPath)
	_ = util.WriteStringToFile("testKey", keyPath)

	// when
	certificate := NewCertificate(domain, certPath, keyPath)
	certStore := NewSecretCertStore(secretManager)
	certStore.Register(certificate)

	// then
	expectedCreateSecretCmdLine :=
		"kubectl -n default create secret tls test-domain --key=/tmp/tls.key --cert=/tmp/tls.crt"
	actualCreateSecretCmdLine := cmdLine.CmdLines[0]

	createdSecretName := kubernetesRunner.FindSecret(domain, "default")
	assert.Contains(t, createdSecretName, domain)
	assert.Equal(t, expectedCreateSecretCmdLine, actualCreateSecretCmdLine)
}


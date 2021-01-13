package secret

import (
	"fmt"
	"os"
	"testing"

	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	// setup
	cmdLine := testutil.FakeKubernetesRunner("test-output")
	kubernetesRunner := cluster.NewKubernetesRunner(cmdLine)
	secretManager := NewKubernetesSecretManager("default", kubernetesRunner)

	// given
	filePath := "/tmp/filepath"
	secretName := "test-secret-name"
	os.Create(filePath)

	// when
	secretManager.CreateFileSecret(secretName, filePath)

	// then
	//TODO: should validate snake case in the secret name (camelCase failures)
	expectedCreateSecretCmdLine :=
		fmt.Sprintf("kubectl create secret generic %s "+
			"--from-file=%s.json=%s -n default", secretName, secretName, filePath)
	actualCreateSecretCmdLine := cmdLine.CmdLines[0]

	assert.Equal(t, expectedCreateSecretCmdLine, actualCreateSecretCmdLine)
}

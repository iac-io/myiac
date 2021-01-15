package secret

import (
	"fmt"
	"testing"

	"github.com/iac-io/myiac/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateFileSecret(t *testing.T) {
	// setup
	cmdLine := testutil.FakeCommandRunner("secret default/test-secret created")
	kubernetesRunner := NewKubernetesRunner(cmdLine)

	// given
	filePath := "/tmp/filepath"
	secretName := "test-Secret-Name"

	// when
	kubernetesRunner.CreateFileSecret(secretName, "default", filePath)

	// then
	expectedDeleteSecretCmdLine := "kubectl delete secret test-Secret-Name -n default"
	expectedCreateSecretCmdLine := fmt.Sprintf("kubectl create secret generic %s "+
		"--from-file=%s.json=%s -n default", secretName, secretName, filePath)
	actualDeleteSecretCmdLine := cmdLine.CmdLines[0]
	actualCreateSecretCmdLine := cmdLine.CmdLines[1]

	assert.Equal(t, expectedCreateSecretCmdLine, actualCreateSecretCmdLine)
	assert.Equal(t, expectedDeleteSecretCmdLine, actualDeleteSecretCmdLine)
}

func TestCreateLiteralSecret(t *testing.T) {
	// setup
	cmdLine := testutil.FakeCommandRunner("output")
	kubernetesRunner := NewKubernetesRunner(cmdLine)

	// given
	secretName := "key-value"
	key := "testkey"
	value := "testValue"
	literalMap := make(map[string]string)
	literalMap[key] = value

	// when
	kubernetesRunner.CreateLiteralSecret(secretName, "default", literalMap)

	// then
	expectedDeleteSecretCmdLine := "kubectl delete secret key-value -n default"
	expectedCreateSecretCmdLine := fmt.Sprintf("kubectl create secret generic %s "+
		"--from-literal=%s=%s -n default", secretName, key, value)
	actualDeleteSecretCmdLine := cmdLine.CmdLines[0]
	actualCreateSecretCmdLine := cmdLine.CmdLines[1]

	assert.Equal(t, expectedCreateSecretCmdLine, actualCreateSecretCmdLine)
	assert.Equal(t, expectedDeleteSecretCmdLine, actualDeleteSecretCmdLine)
}

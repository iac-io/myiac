package secret

import (
	"fmt"
	"log"
	"strings"

	"github.com/iac-io/myiac/internal/commandline"
)

type KubernetesSecretRunner interface {
	CreateTlsSecret(name string, namespace string, keyFile string, certFile string)
	CreateFileSecret(name string, namespace string, filePath string)
	CreateLiteralSecret(name string, namespace string, literalsMap map[string]string)
	FindSecret(name string, namespace string) string
}

type kubernetesRunner struct {
	cmdRunner commandline.CommandRunner
}

// NewKubernetesRunner creates a command runner for kubernetes secret-related operations
func NewKubernetesRunner(commandRunner commandline.CommandRunner) KubernetesSecretRunner {
	return &kubernetesRunner{cmdRunner: commandRunner}
}

// CreateTlsSecret create a TLS secret in Kubernetes, used to store SSL certificates from its cert and key files
// note: it deletes any existing secret with the same name in the same namespace
func (kr kubernetesRunner) CreateTlsSecret(name string, namespace string, keyFile string, certFile string) {
	deleteSecret(kr.cmdRunner, name, namespace)
	keysArg := ""

	fmt.Printf("Adding key file: %s -> %s", keyFile, "*****\n")
	keyArg := fmt.Sprintf("--key=%s", keyFile)
	fmt.Printf("Adding cert file: %s -> %s", certFile, "*****\n")
	certArg := fmt.Sprintf("--cert=%s", certFile)

	keysArg = strings.TrimSpace(keysArg)
	argsArray := []string{"-n", namespace, "create", "secret", "tls", name, keyArg, certArg}

	kr.cmdRunner.SetupWithoutOutput("kubectl", argsArray)
	kr.cmdRunner.Run()
}

// CreateSecret creates a generic secret in a Kubernetes namespace from an existing
// service account key JSON file
//
// See: https://stackoverflow.com/questions/45879498/how-can-i-update-a-secret-on-kubernetes-when-it-is-generated-from-a-file
// Example: kubectl create secret generic firestore-key --from-file=key.json=/path/to/moneycol-firestore-collections-api.json
func (kr kubernetesRunner) CreateFileSecret(name string, namespace string, jsonKeyPath string) {
	deleteSecret(kr.cmdRunner, name, namespace)
	fromFileArg := fmt.Sprintf("--from-file=%s.json=%s", name, jsonKeyPath)
	cmdLine := fmt.Sprintf("kubectl create secret generic %s %s -n %s", name, fromFileArg, namespace)
	kr.cmdRunner.SetupCmdLine(cmdLine)
	kr.cmdRunner.Run()
}

func (kr kubernetesRunner) FindSecret(name string, namespace string) string {
	argsArray := []string{"get", "secret", name, "-n", namespace}
	kr.cmdRunner.SetupWithoutOutput("kubectl", argsArray)
	cmdOutput := kr.cmdRunner.Run()
	return cmdOutput.Output
}

// kubectl create secret generic dev-db-secret --from-literal=username=devuser --from-literal=password='S!B\*d$zDsb='
func (kr kubernetesRunner) CreateLiteralSecret(name string, namespace string, literalsMap map[string]string) {
	log.Printf("kubernetes secret: clearing any existing secret with name %s in namespace %s first...", name, namespace)
	deleteSecret(kr.cmdRunner, name, namespace)
	fromLiteralArg := ""
	for k, v := range literalsMap {
		fmt.Printf("Adding secret literal: %s -> %s", k, "*****\n")
		fromLiteralArg += fmt.Sprintf("--from-literal=%s=%s ", k, v)
	}

	fromLiteralArg = strings.TrimSpace(fromLiteralArg)

	cmdLine := fmt.Sprintf("kubectl create secret generic %s %s -n %s", name, fromLiteralArg, namespace)
	kr.cmdRunner.SetupCmdLine(cmdLine)
	kr.cmdRunner.SetSuppressOutput(true)
	kr.cmdRunner.Run()
}

func deleteSecret(cmdRunner commandline.CommandRunner, name string, namespace string) {
	cmdLine := fmt.Sprintf("kubectl delete secret %s -n %s", name, namespace)
	cmdRunner.SetupCmdLine(cmdLine)
	cmdRunner.SetSuppressOutput(true)
	cmdRunner.IgnoreError(true)
	cmdRunner.Run()
}

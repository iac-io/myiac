package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/commandline"
	"github.com/dfernandezm/myiac/internal/util"
	"strings"
)

func GetInternalIpsForNodes() []string {
	json := executeGetIpsCmd()
	ips := getAllIps(json, true)
	fmt.Printf("Internal IPs for nodes in cluster are: %v\n", ips)
	return ips
}

func GetAllPublicIps() []string {
	json := executeGetIpsCmd()
	ips := getAllIps(json, false)
	fmt.Printf("Public IPs for nodes in cluster are: %v\n", ips)
	return ips
}

func GetPods() {
	baseArgs := "get pods"
	var argsArray []string = strings.Fields(baseArgs)
	cmd := commandline.New("kubectl", argsArray)
	cmd.Run()
}

func executeGetIpsCmd() map[string]interface{} {
	argsArray := []string{"get", "nodes", "-o", "json"}
	cmd := commandline.New("kubectl", argsArray)
	cmd.IsSuppressOutput = true
	cmdResult := cmd.Run()
	cmdOutput := cmdResult.Output
	json := util.Parse(cmdOutput)
	return json
}

func getAllIps(json map[string]interface{}, internal bool) []string {
	indexOfAddress := 1
	if internal {
		indexOfAddress = 0
	} 
	allNodes := util.GetJsonArray(json, "items")
	var ips []string
	for _, node := range allNodes {
		status := util.GetJsonObject(node, "status")
		addresses := util.GetJsonArray(status, "addresses")
		ip := util.GetStringValue(addresses[indexOfAddress], "address") 
		ips = append(ips, ip)
	}
	return ips
}

// kubectl create secret generic dev-db-secret --from-literal=username=devuser --from-literal=password='S!B\*d$zDsb='
func CreateSecretFromLiteral(name string, namespace string, literals map[string]string) {
	deleteSecret(name, namespace)
	fromLiteralArg := ""
	for k, v := range literals {
		fmt.Printf("Adding secret literal: %s -> %s", k, "*****\n")
		fromLiteralArg += fmt.Sprintf("--from-literal=%s=%s ", k, v)
	}

	fromLiteralArg = strings.TrimSpace(fromLiteralArg)
	argsArray := []string{"create", "secret", "generic", name, fromLiteralArg, "-n", namespace}
	cmd := commandline.New("kubectl", argsArray)
	cmd.IsSuppressOutput = true
	cmd.Run()
}

// ================= Move this to new file ======================

type KubernetesRunner interface {
	CreateTlsSecret(name string, namespace string, keyFile string, certFile string)
	CreateFileSecret(name string, namespace string, filePath string)
	FindSecret(name string, namespace string) string
}

type kubernetesRunner struct {
	cmdRunner commandline.CommandRunner
}

func NewKubernetesRunner(commandRunner commandline.CommandRunner) *kubernetesRunner {
	return &kubernetesRunner{cmdRunner:commandRunner}
}

// CreateTlsSecret create a TLS based secret in Kubernetes, used to store SSL certificates from its cert and key files
func (kr kubernetesRunner) CreateTlsSecret(name string, namespace string, keyFile string, certFile string) {
	deleteSecret(name, namespace)
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
	deleteSecret(name, namespace)
	fromFileArg := fmt.Sprintf("--from-file=%s.json=%s", name, jsonKeyPath)
	cmdLine := fmt.Sprintf("kubectl create secret generic %s %s -n %s", name, fromFileArg, namespace)
	kr.cmdRunner.SetupCmdLine(cmdLine)
	kr.cmdRunner.Run()
}

func (kr kubernetesRunner) FindSecret(name string, namespace string) string {
	argsArray := []string{"get", "secret", name,"-n", namespace}
	kr.cmdRunner.SetupWithoutOutput("kubectl", argsArray)
	cmdOutput := kr.cmdRunner.Run()
	return cmdOutput.Output
}

func deleteSecret(name string, namespace string) {
	argsArray := []string{"delete", "secret", name, "-n", namespace}
	cmd := commandline.New("kubectl", argsArray)
	cmd.IsSuppressOutput = true
	cmd.IgnoreError(true)
	cmd.Run()
}



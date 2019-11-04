package main

//% go build -o $GOPATH/bin/myiac github.com/dfernandezm/myiac/app
import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strings"
)

//https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
//https://golang.org/doc/code.html
func main() {
	fmt.Printf("MyIaC - Infrastructure as Code\n")
	runtime := RuntimeProperties{}
	setupEnvironment()
	configureDocker()
	tagDockerImage(&runtime)
	pushDockerImage(&runtime)
	//setupKubernetes()
	//getPods()
	//labelElasticsearchNodes()
	//labelDockerImage()
}

func setupEnvironment() {
	// create a service account and download it:
	keyLocation := getHomeDir() + "/account.json"
	baseArgs := "auth activate-service-account --key-file %s"
	baseArgsTmpl := fmt.Sprintf(baseArgs, keyLocation)
	var argsArray []string = strings.Fields(baseArgsTmpl)
	command("gcloud", argsArray)
}

func getPods() {
	baseArgs := "get pods"
	var argsArray []string = strings.Fields(baseArgs)
	command("kubectl", argsArray)
}

func setupKubernetes() {
	project := "moneycol"
	zone := "europe-west1-b"
	clusterName := "moneycol-main"
	//split -- needs to be an array
	clusterCredentialsPart := "container clusters get-credentials"
	argsStr := fmt.Sprintf("%s %s --zone %s --project %s", clusterCredentialsPart, clusterName, zone, project)
	argsArray := strings.Fields(argsStr)
	command("gcloud", argsArray)
}

func labelElasticsearchNodes() {
	nodeName := "gke-moneycol-main-elasticsearch-pool-b8711571-k359"
	label := "type=elasticsearch"
	labelCmdTpl := "label nodes %s %s"
	argsStr := fmt.Sprintf(labelCmdTpl, nodeName, label)
	argsArray := strings.Fields(argsStr)
	command("kubectl", argsArray)
}

func configureDocker() {
	action := "auth configure-docker"
	cmdTpl := "%s"
	argsStr := fmt.Sprintf(cmdTpl, action)
	argsArray := strings.Fields(argsStr)
	command("gcloud", argsArray)
}

func tagDockerImage(runtime *RuntimeProperties) {
	imageToTag := "c58e8ce1a62e"
	projectId := "moneycol"
	projectRepository := "gcr.io"
	containerName := "moneycol-server"
	tag := "0.1.0-alpha1"
	containerFullName := fmt.Sprintf("%s/%s/%s:%s", projectRepository, projectId, containerName, tag)
	fmt.Printf("The image tag to push is: %s\n", containerFullName)
	dockerTagCmdPart := "tag"
	argsStr := fmt.Sprintf("%s %s %s\n", dockerTagCmdPart, imageToTag, containerFullName)
	argsArray := strings.Fields(argsStr)
	fmt.Printf("Array of args %s\n", argsArray)
	err := command("docker", argsArray)
	if err != nil {
		log.Fatalf("Command '%s' failed with error %s\n", "docker "+argsStr, err)
	}
	runtime.SetDockerImage(containerFullName)
	fmt.Printf("Docker image has been tagged with: %s\n", runtime.GetDockerImage())
}

func pushDockerImage(runtime *RuntimeProperties) {
	fmt.Printf("Pushing previously built docker image: %s", runtime.GetDockerImage())
	argsStr := fmt.Sprintf("%s %s", "push", runtime.GetDockerImage())
	argsArray := strings.Fields(argsStr)
	command("docker", argsArray)
}

func command(command string, arguments []string) error {
	cmd := exec.Command(command, arguments...)
	//cmd.Dir = "/Users/david"
	fmt.Printf("Executing [ %s ]\n", string(strings.Join(cmd.Args, " ")))
	out, err := cmd.CombinedOutput() //TODO: get stderr and stdout in separate strings
	if err != nil {
		fmt.Printf("Output: \n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	fmt.Printf("Output: \n%s\n", string(out))
	return nil
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// ------------------- separate this -----

type RuntimeProperties struct {
	DockerImage string
}

func (rp *RuntimeProperties) SetDockerImage(dockerImage string) {
	rp.DockerImage = dockerImage
}

func (rp RuntimeProperties) GetDockerImage() string {
	return rp.DockerImage
}

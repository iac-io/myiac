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
	setupEnvironment()
	setupKubernetes()
	getPods()
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

func command(command string, arguments []string) {
	cmd := exec.Command(command, arguments...)
	//cmd.Dir = "/Users/david"
	fmt.Printf("Executing [ %s ]\n", string(strings.Join(cmd.Args, " ")))
	out, err := cmd.CombinedOutput() //TODO: get stderr and stdout in separate strings
	if err != nil {
		fmt.Printf("Output: \n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("Output: \n%s\n", string(out))
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

package main

//% go build -o $GOPATH/bin/myiac github.com/dfernandezm/myiac/app
import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

//https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func main() {
	fmt.Printf("MyIaC - Infrastructure as Code\n")
	//command("ls", "-lah")
	setupEnvironment()
	//https://golang.org/doc/code.html
}

func setupEnvironment() {
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
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Output: \n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("Output: \n%s\n", string(out))
}

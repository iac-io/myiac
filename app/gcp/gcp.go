package gcp

import (
	"fmt"

	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

func SetupEnvironment() {
	keyLocation := util.GetHomeDir() + "/account.json"
	baseArgs := "auth activate-service-account --key-file %s"
	var argsArray []string = util.StringTemplateToArgsArray(baseArgs, keyLocation)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
}

func ConfigureDocker() {
	action := "auth configure-docker"
	cmdTpl := "%s"
	argsArray := util.StringTemplateToArgsArray(cmdTpl, action)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
}

func SetupKubernetes(project string, zone string, environment string) {
	action := "container clusters get-credentials"
	// gcloud container clusters get-credentials [cluster-name]
	cmdTpl := "%s %s --zone %s --project %s"
	clusterName := fmt.Sprintf("%s-%s", project, environment)

	argsArray := util.StringTemplateToArgsArray(cmdTpl, action, clusterName, zone, project)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
	fmt.Println("Kubernetes setup completed")
}

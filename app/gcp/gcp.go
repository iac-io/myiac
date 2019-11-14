package gcp

import (
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

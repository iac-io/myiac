package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/gcp"
	"github.com/dfernandezm/myiac/app/util"
	"time"
)

func InstallHelm() {
	fmt.Println("Installing Helm in newly created cluster")

	currentDir := util.CurrentExecutableDir()
	helperScriptsLocation := currentDir + "/helperScripts"
	fmt.Printf("Helper scripts path is %s", helperScriptsLocation)

	action := "./install-helm.sh"
	cmd := commandline.NewWithWorkingDir(action, []string{}, helperScriptsLocation)
	cmd.Run()
}

func InitTerraform() {
	currentDir := util.CurrentExecutableDir()
	initScriptLocation := currentDir + "/terraform"
	initScript := "./init.sh"

	cmd := commandline.NewWithWorkingDir(initScript, []string{}, initScriptLocation)
	cmd.Run()
}

//TODO: pass variables into the TF template/clustervars
func CreateCluster(project string, env string, zone string) {
	currentDir := util.CurrentExecutableDir()
	tfFileLocation := currentDir + "/terraform/cluster"

	varFileLoc := tfFileLocation + "/cluster.tfvars"
	argsArray := util.StringTemplateToArgsArray("%s %s", "plan", "-var-file="+varFileLoc)
	cmd := commandline.NewWithWorkingDir("terraform", argsArray, tfFileLocation)
	cmd.Run()

	argsArray = util.StringTemplateToArgsArray("%s %s %s", "apply", "-var-file="+varFileLoc, "-auto-approve")
	cmd = commandline.NewWithWorkingDir("terraform", argsArray, tfFileLocation)
	cmd.Run()

	fmt.Printf("Kubernetes cluster created through Terraform from %s\n", tfFileLocation)

	fmt.Println("Installing Helm into newly created cluster...")

	fmt.Println("Waiting 10 seconds for cluster to stabilize before installing Helm")
	time.Sleep(10 * time.Second)

	fmt.Println("Setting up newly created Kubernetes cluster")
	gcp.SetupKubernetes(project, env, zone)

	InstallHelm()
}

func DestroyCluster() {
	currentDir := util.CurrentExecutableDir()

	tfFileLocation := currentDir + "/terraform/cluster"
	varFileLoc := tfFileLocation + "/cluster.tfvars"

	fmt.Println("Waiting 5 seconds before destroying cluster...")
	time.Sleep(5 * time.Second)

	argsArray := util.StringTemplateToArgsArray("%s %s %s", "destroy", "-var-file="+varFileLoc, "-auto-approve")
	cmd := commandline.NewWithWorkingDir("terraform", argsArray, tfFileLocation)
	cmd.Run()

	fmt.Println("Kubernetes cluster deleted through Terraform")
}

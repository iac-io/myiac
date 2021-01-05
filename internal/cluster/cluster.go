package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/commandline"
	"github.com/dfernandezm/myiac/internal/gcp"
	"github.com/dfernandezm/myiac/internal/util"
	"log"
	"os"
	"time"
	)

var (
	tfvarsPath = util.CurrentExecutableDir() + "/internal/terraform/cluster"
	tfvarsFile = tfvarsPath + "/cluster.tfvars"
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

func InitTerraform(tf string, project string, env string) {
	// Create Bucket if not present in the system
	err := gcp.CreateGCSBucket(project, env)
	if err != nil {
		log.Fatal(err)
	}
	// Check if terraform initialized
	if _, err := os.Stat(tf+"/.terraform"); os.IsNotExist(err) {
		fmt.Printf("Terraform not Initialized: in %v: Intializong now...", tf+"/.terraform\n")
		argsArray := util.StringTemplateToArgsArray("%s", "init")
		cmd := commandline.NewWithWorkingDir("terraform", argsArray, tf)
		cmd.Run()
	}
}

func PlanTerraform(tp string, tf string)  {
	argsArray := util.StringTemplateToArgsArray("%s %s", "plan", "-var-file="+tf)
	cmd := commandline.NewWithWorkingDir("terraform", argsArray, tp)
	cmd.Run()
}

func ApplyTerraform(tp string, tf string)  {
	argsArray := util.StringTemplateToArgsArray("%s %s %s", "apply", "-var-file="+tf, "-auto-approve")
	cmd := commandline.NewWithWorkingDir("terraform", argsArray, tp)
	cmd.Run()
}

//TODO: pass variables into the TF template/clustervars
func CreateCluster(project string, env string, zone string, flag bool) {
	InitTerraform(tfvarsPath, project, env)
	if flag {
		fmt.Println("Running Plan only due to --noop option")
		PlanTerraform(tfvarsPath,tfvarsFile)
		os.Exit(0)
	} else {
		ApplyTerraform(tfvarsPath,tfvarsFile)
	}

	fmt.Printf("Kubernetes cluster created through Terraform from %s\n", tfvarsPath)

	fmt.Println("Installing Helm into newly created cluster...")

	fmt.Println("Waiting 10 seconds for cluster to stabilize before installing Helm")
	time.Sleep(10 * time.Second)

	fmt.Println("Setting up newly created Kubernetes cluster")
	//gcp.SetupKubernetes(project, env, zone)

	InstallHelm()
}

func DestroyCluster(project string, env string, zone string) {
	InitTerraform(tfvarsPath, project, env)
	fmt.Println("Waiting 5 seconds before destroying cluster...")
	time.Sleep(5 * time.Second)
	argsArray := util.StringTemplateToArgsArray("%s %s %s", "destroy", "-var-file="+tfvarsFile, "-auto-approve")
	cmd := commandline.NewWithWorkingDir("terraform", argsArray, tfvarsPath)
	cmd.Run()
	fmt.Println("Kubernetes cluster deleted through Terraform")
	//TODO: Need to have option (flag) to be able to delete bucket after cluster is destroyed
	//Code below is ready but GCS cannot remove not empty buckets
	//err := gcp.DeleteGCSBucket(project, env)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

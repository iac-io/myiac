package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strings"
	"github.com/dfernandezm/myiac/app/docker"
	"github.com/dfernandezm/myiac/app/gcp"
	props "github.com/dfernandezm/myiac/app/properties"
)

const (
	projectId            = "moneycol"
	projectRepositoryUrl = "gcr.io"
)

func main() {
	fmt.Printf("MyIaC - Infrastructure as Code\n")
	runtime := props.NewRuntime()
	dockerProps := props.DockerProperties{ProjectRepoUrl: projectRepositoryUrl, ProjectId: projectId}

	gcp.SetupEnvironment()
	gcp.ConfigureDocker()

	imageId := "af898a99ee67"
	commitHash := "bb9c6a4"
	appName := "moneycol-frontend"
	version := "0.1.0"

	docker.TagImage(&runtime, &dockerProps, imageId, commitHash, appName, version)
	docker.PushImage(&runtime)

	//setupKubernetes()
	//getPods()
	//labelNodes("elasticsearch")
	//labelNodes("applications")
	//labelDockerImage()

	// --- MoneyCol server ---
	// basePath := getHomeDir() + "/development/repos/moneycol/server/deploy"
	// appName := "moneycol-server"
	// chartPath := fmt.Sprintf("%s/%s/chart", basePath, appName)
	// moneyColServerDeploy := Deployment{AppName: appName, ChartPath: chartPath, DryRun: false}
	// deployApp(&moneyColServerDeploy)

	// --- Traefik ---
	//chartPath = "stable/traefik"
	//--set dashboard.enabled=true,dashboard.domain=dashboard.localhost
	//traefikDeploy := Deployment{AppName: appName, ChartPath: chartPath, DryRun: false}
	//traefikDeploy.HelmSetParams = "dashboard.enabled=true,dashboard.domain=dashboard.localhost"

	//deployApp(&traefikDeploy)
}

// commit hash: git rev-parse HEAD | cut -c1-7
// func tagDockerImageNew(runtimeProperties *props.RuntimeProperties, imageId string, commitHash string, imageRepo string, appVersion string) {
// 	img := docker.NewImageDefinition(projectRepositoryUrl, projectId, imageId, commitHash, imageRepo, appVersion)
// 	img.TagImage(runtimeProperties)
// }

// func pushDockerImageNew(runtimeProperties *props.RuntimeProperties) {
// 	docker.PushImage()
// }

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

func labelNodes(nodeType string) {
	//slice vs array: https://blog.golang.org/go-slices-usage-and-internals
	var nodeNames []string
	var label string
	nodeNamesEs := []string{"gke-moneycol-main-elasticsearch-pool-b8711571-k359"}
	nodeNamesApps := []string{"gke-moneycol-main-main-pool-ac0c4442-57ff",
		"gke-moneycol-main-main-pool-ac0c4442-pq57",
		"gke-moneycol-main-main-pool-ac0c4442-q1t7"}

	if nodeType == "elasticsearch" {
		nodeNames = nodeNamesEs
		label = "type=elasticsearch"
	} else if nodeType == "applications" {
		nodeNames = nodeNamesApps
		label = "type=applications"
	}

	labelCmdTpl := "label nodes %s %s --overwrite\n"

	//note: range (like everything in go) copies by value the slice
	for _, nodeName := range nodeNames {
		argsStr := fmt.Sprintf(labelCmdTpl, nodeName, label)
		fmt.Printf("Labelling args: %s", argsStr)
		argsArray := strings.Fields(argsStr)
		command("kubectl", argsArray)
	}
}

func pushDockerImage(runtime *RuntimeProperties) {
	fmt.Printf("Pushing previously built docker image: %s\n", runtime.GetDockerImage())
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

func deployApp(deployment *Deployment) {
	cmdExec := "helm"

	//argsTpl := "install %s"
	argsTpl := "upgrade solemn-fly %s"

	argsStr := fmt.Sprintf(argsTpl, deployment.ChartPath)

	if len(deployment.HelmSetParams) != 0 {
		argsStr = fmt.Sprintf("%s --set %s", argsStr, deployment.HelmSetParams)
	}

	argsArray := strings.Fields(argsStr)
	command(cmdExec, argsArray)
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

type Deployment struct {
	AppName   string
	ChartPath string
	DryRun    bool
	HelmSetParams string
}

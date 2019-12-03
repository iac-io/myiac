package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strings"

	"github.com/dfernandezm/myiac/app/cluster"
	"github.com/dfernandezm/myiac/app/docker"
	"github.com/dfernandezm/myiac/app/gcp"
	props "github.com/dfernandezm/myiac/app/properties"
)

const (
	project              = "moneycol"
	clusterZone          = "europe-west1-b"
	environment          = "dev"
	projectId            = "moneycol"
	projectRepositoryUrl = "gcr.io"
)

func main() {
	fmt.Printf("MyIaC - Infrastructure as Code\n")
	gcp.SetupEnvironment()
	gcp.ConfigureDocker()
	gcp.SetupKubernetes("moneycol", "europe-west1-b", "dev")
	//cluster.GetInternalIpsForNodes()
	//cluster.InstallHelm()
	deployApps()

	// ------ Docker workflows  -------
	//runtime := props.NewRuntime()
	//dockerWorkflows(&runtime)

	// --- Various kubernetes setups ---
	//setupKubernetes()
	//getPods()
	//labelNodes("elasticsearch")
	//labelNodes("applications")
	//labelDockerImage()
}

// commit hash: git rev-parse HEAD | cut -c1-7
// func tagDockerImageNew(runtimeProperties *props.RuntimeProperties, imageId string, commitHash string, imageRepo string, appVersion string) {
// 	img := docker.NewImageDefinition(projectRepositoryUrl, projectId, imageId, commitHash, imageRepo, appVersion)
// 	img.TagImage(runtimeProperties)
// }

// func pushDockerImageNew(runtimeProperties *props.RuntimeProperties) {
// 	docker.PushImage()
// }

func dockerWorkflows(runtime *props.RuntimeProperties) {
	imageId := "b63a014d5aa6"
	commitHash := "9ff9cd6"
	appName := "moneycol-server"
	version := "0.1.0"
	dockerProps := props.DockerProperties{ProjectRepoUrl: projectRepositoryUrl, ProjectId: projectId}

	docker.TagImage(runtime, &dockerProps, imageId, commitHash, appName, version)
	docker.PushImage(runtime)
}

func deployApps() {
	// -- Elasticsearch ---
	//deployElasticsearch()

	// --- MoneyCol server ---
	//deployMoneyColServer()

	// --- MoneyCol frontend ---
	//deployMoneyColFrontend()

	// --- Traefik ---
	deployTraefik()
}

func deployTraefik() {
	releaseName := "opining-frog"
	moneycolPath := "/development/repos/moneycol/"
	basePath := getHomeDir() + moneycolPath + "server/deploy"
	appName := "traefik"
	chartPath := fmt.Sprintf("%s/%s/chart", basePath, appName)

	//TODO: Set paramaters, separate this into helm.go
	helmSetParams := make(map[string]string)
	internalIps := cluster.GetInternalIpsForNodes()
	internalIpsForHelmSet := "{" + strings.Join(internalIps, ",") + "}"
	fmt.Printf("Internal IPs to set for helm are: %s", internalIpsForHelmSet)
	helmSetParams["externalIps"] = "\"" + internalIpsForHelmSet + "\""

	deployment := Deployment{AppName: appName, ChartPath: chartPath, 
							DryRun: false, 
							HelmReleaseName: releaseName,
							HelmSetParams: helmSetParams}
	deployApp(&deployment)
}

func deployMoneyColFrontend() {
	releaseName := "esteemed-peacock"
	moneycolPath := "/development/repos/moneycol/"
	basePath := getHomeDir() + moneycolPath + "frontend/deploy"
	appName := "moneycolfrontend"
	chartPath := fmt.Sprintf("%s/%s/chart", basePath, appName)
	moneyColFrontendDeploy := Deployment{AppName: appName, ChartPath: chartPath, DryRun: false, HelmReleaseName: releaseName}
	deployApp(&moneyColFrontendDeploy)
}

func deployElasticsearch() {
	releaseName := ""
	moneycolPath := "/development/repos/moneycol/"
	basePath := getHomeDir() + moneycolPath + "server/deploy"
	appName := "elasticsearch"
	chartPath := fmt.Sprintf("%s/%s/chart", basePath, appName)
	elasticsearchDeploy := Deployment{AppName: appName, ChartPath: chartPath, DryRun: false, HelmReleaseName: releaseName}
	deployApp(&elasticsearchDeploy)
}

func deployMoneyColServer() {
	releaseName := "ponderous-lion"
	moneycolPath := "/development/repos/moneycol/"
	basePath := getHomeDir() + moneycolPath + "server/deploy"
	appName := "moneycolserver"
	chartPath := fmt.Sprintf("%s/%s/chart", basePath, appName)
	moneyColServerDeploy := Deployment{AppName: appName, ChartPath: chartPath, DryRun: false, HelmReleaseName: releaseName}
	deployApp(&moneyColServerDeploy)
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

	var argsTpl = ""
	if deployment.HelmReleaseName == "" {
		argsTpl = "install %s"
	} else {
		argsTpl = "upgrade " + deployment.HelmReleaseName + " %s"
	}

	argsStr := fmt.Sprintf(argsTpl, deployment.ChartPath)

	if len(deployment.HelmSetParams) != 0 {
		setParams := ""
		for k, v := range deployment.HelmSetParams {
			setParams += setParams + "--set " + k + "=" + v + " " 
		}

		argsStr = fmt.Sprintf("%s %s", argsStr, setParams)
	}

	if (deployment.DryRun) {
		argsStr += " --debug --dry-run"
	}

	argsArray := strings.Fields(argsStr)
	command(cmdExec, argsArray)
}

// ------------------- separate this -----

type Deployment struct {
	AppName         string
	ChartPath       string
	DryRun          bool
	HelmSetParams   map[string]string // key value pairs, get its own struct soon
	HelmReleaseName string
}

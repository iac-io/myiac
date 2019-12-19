package deploy

import (
	"fmt"
	"strings"
	"os"
	"github.com/dfernandezm/myiac/app/gcp"
	"github.com/dfernandezm/myiac/app/cluster"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

type Deployment struct {
	AppName          string
	ChartPath        string
	DryRun           bool
	HelmSetParams    map[string]string // key value pairs
	HelmReleaseName  string
	HelmValuesParams []string // yaml filenames to pass as --values
}

// DeployApp deploys the given appName into the given environment
// A helm chart folder named appName must exist in charts base path
func DeployApp(appName string, environment string) {
	//TODO: generify
	if appName == "all" {
		deployApps(environment)
	}

	if appName == "moneycolserver" {
		deployMoneyColServer()
	}

	if appName == "moneycolfrontend" {
		deployMoneyColFrontend()
	}

	if appName == "elasticsearch" {
		deployElasticsearch()
	}

	if appName == "traefik" {
		deployTraefik(environment)
	}

	if appName == "traefik-dev" {
		deployTraefikDev()
	}
}

func deployApps(environment string) {
	deployElasticsearch()
	deployMoneyColServer()
	deployMoneyColFrontend()
	deployTraefik(environment)
}

func deployMoneyColFrontend() {
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(cmdRunner)
	releaseName := helmDeployer.ReleaseFor("moneycolfrontend")
	
	appName := "moneycolfrontend"
	baseChartsPath := getBaseChartsPath()
	chartPath := fmt.Sprintf("%s/%s", baseChartsPath, appName)
	moneyColFrontendDeploy := Deployment{AppName: appName, ChartPath: chartPath, 
								DryRun: false, HelmReleaseName: releaseName}
	deployApp(&moneyColFrontendDeploy)
}

func deployElasticsearch() {
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(cmdRunner)
	releaseName := helmDeployer.ReleaseFor("elasticsearch")
	
	appName := "elasticsearch"
	baseChartsPath := getBaseChartsPath()
	chartPath := fmt.Sprintf("%s/%s", baseChartsPath, appName)
	
	elasticsearchDeploy := Deployment{AppName: appName, ChartPath: chartPath, 
										DryRun: false, HelmReleaseName: releaseName}
	deployApp(&elasticsearchDeploy)
}

func deployMoneyColServer() {
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(cmdRunner)
	releaseName := helmDeployer.ReleaseFor("moneycolserver")
	appName := "moneycolserver"

	baseChartsPath := getBaseChartsPath()
	chartPath := fmt.Sprintf("%s/%s", baseChartsPath, appName)

	moneyColServerDeploy := Deployment{AppName: appName, ChartPath: chartPath, 
										DryRun: false, HelmReleaseName: releaseName}
	deployApp(&moneyColServerDeploy)
}

func deployTraefik(environment string) {
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(cmdRunner)
	releaseName := helmDeployer.ReleaseFor("traefik")
	appName := "traefik"

	baseChartsPath := getBaseChartsPath()
	chartPath := fmt.Sprintf("%s/%s", baseChartsPath, appName)
	
	//TODO: Set paramaters, separate this
	helmSetParams := make(map[string]string)
	internalIps := cluster.GetInternalIpsForNodes()

	// very flaky --set for ips like this: --set externalIps={ip1\,ip2\,ip3}
	internalIpsForHelmSet := "{" + strings.Join(internalIps, "\\,") + "}"
	helmSetParams["externalIps"] = internalIpsForHelmSet
	deployment := Deployment{AppName: appName, ChartPath: chartPath,
		DryRun:          false,
		HelmReleaseName: releaseName,
		HelmSetParams:   helmSetParams}

	deployApp(&deployment)
	
	if (environment == "dev") {
		deployTraefikDev()

		// once deployed, repoint dev DNS to any public IP of nodes
		changeDevDNS()
	}
}

func deployTraefikDev() {
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(cmdRunner)
	releaseName := helmDeployer.ReleaseFor("traefik-dev")
	appName := "traefik-dev"
	baseChartsPath := getBaseChartsPath()
	chartPath := fmt.Sprintf("%s/%s", baseChartsPath, appName)

	helmSetParams := make(map[string]string)
	internalIps := cluster.GetInternalIpsForNodes()

	// very flaky --set for ips like this: --set externalIps={ip1\,ip2\,ip3}
	internalIpsForHelmSet := "{" + strings.Join(internalIps, "\\,") + "}"
	helmSetParams["externalIps"] = internalIpsForHelmSet
	deployment := Deployment{AppName: appName, ChartPath: chartPath,
		DryRun:          false,
		HelmReleaseName: releaseName,
		HelmSetParams:   helmSetParams}

	deployApp(&deployment)
	changeDevDNS()
}

func getBaseChartsPath() string {
	chartsPath := os.Getenv("CHARTS_PATH")
	if chartsPath != "" {
		return chartsPath
	}

	chartsPath = util.CurrentExecutableDir() + "/charts"
	return chartsPath
}

func changeDevDNS() {
	publicIps := cluster.GetAllPublicIps()
	aPublicIP := publicIps[0] // any public ip works for this
	applyDNSThroughSdk(aPublicIP)
}

func applyDNSThroughSdk(newIP string) {
	fmt.Printf("Applying changes to DNS for new IP: %s", newIP)
	dnsService := gcp.NewGoogleCloudDNSService("moneycol","money-zone-free")
	dnsService.UpsertDNSRecord("A", "dev.moneycol.ml", newIP)
}

func applyDNSUsingTerraform(newIP string) {
	moneycolPath := "/development/repos/moneycol/"
	deployPath := util.GetHomeDir() + moneycolPath + "server/deploy"
	devDNSTfFile := deployPath + "/terraform/dns"
	cluster.ApplyDnsIpChange(devDNSTfFile, newIP)
}

func deployApp(deployment *Deployment) {

	var argsTpl = ""
	if deployment.HelmReleaseName == "" {
		argsTpl = "install %s"
	} else {
		argsTpl = "upgrade " + deployment.HelmReleaseName + " %s"
	}

	argsStr := fmt.Sprintf(argsTpl, deployment.ChartPath)

	if len(deployment.HelmValuesParams) > 0 {
		valuesParams := ""
		for _, filePath := range deployment.HelmValuesParams {
			valuesParams += valuesParams + "--values " + filePath + " "
		}
		valuesParams = strings.TrimSpace(valuesParams)
		argsStr = fmt.Sprintf("%s %s", argsStr, valuesParams)
	}

	if len(deployment.HelmSetParams) > 0 {
		setParams := ""
		for k, v := range deployment.HelmSetParams {
			setParams += setParams + "--set " + k + "=" + v + " "
		}
		setParams = strings.TrimSpace(setParams)
		argsStr = fmt.Sprintf("%s %s", argsStr, setParams)
	}

	if deployment.DryRun {
		argsStr += " --debug --dry-run"
	}

	argsArray := strings.Fields(argsStr)
	cmd := commandline.New("helm", argsArray)
	cmd.Run()
	fmt.Printf("Finished deploying app: %s\n\n", deployment.AppName)
}

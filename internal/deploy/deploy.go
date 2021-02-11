package deploy

import (
	"fmt"
	"github.com/iac-io/myiac/internal/util"
	"os"
	"strings"

	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/gcp"
)

type Deployment struct {
	AppName          string
	ChartPath        string
	DryRun           bool
	HelmSetParams    map[string]string // key value pairs
	HelmReleaseName  string
	HelmValuesParams []string // yaml filenames to pass as --values
}

func deployApps(environment string) {
	//TODO: Read apps from manifest
}

func getNodesInternalIpsAsHelmParams() map[string]string {
	helmSetParams := make(map[string]string)
	internalIps := cluster.GetInternalIpsForNodes()

	// very flaky --set for ips like this: --set externalIps={ip1\,ip2\,ip3}
	internalIpsForHelmSet := "{" + strings.Join(internalIps, "\\,") + "}"
	helmSetParams["externalIps"] = internalIpsForHelmSet
	return helmSetParams
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
	aPublicIP := publicIps[0] // any public ip works for this as it's clusterIP
	applyDNSThroughSdk(aPublicIP)
}

func applyDNSThroughSdk(newIP string) {
	fmt.Printf("Applying changes to DNS for new IP: %s", newIP)
	dnsService := gcp.NewGoogleCloudDNSService("moneycol", "money-zone-free")
	dnsService.UpsertDNSRecord("A", "dev.moneycol.ml", newIP)
}

// moneycolfrontend, moneycolserver, elasticsearch, traefik, traefik-dev, collections-api
func Deploy(appName string, environment string, propertiesMap map[string]string, dryRun bool) {
	helmSetParams := make(map[string]string)
	if appName == "traefik" || appName == "traefik-dev" || appName == "test" {
		helmSetParams = getNodesInternalIpsAsHelmParams()
	}

	//TODO: Add properties to helmSetParams or values
	addPropertiesToSetParams(helmSetParams, propertiesMap)
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(getBaseChartsPath(), cmdRunner)
	deployment := HelmDeployment{
		AppName:       appName,
		Environment:   environment,
		HelmSetParams: helmSetParams,
		DryRun:        dryRun,
	}
	helmDeployer.Deploy(&deployment)

	//if appName == "traefik" || appName == "traefik-dev" {
	//	changeDevDNS()
	//}
}

func addPropertiesToSetParams(helmSetParams map[string]string, propertiesMap map[string]string) {
	for k, v := range propertiesMap {
		fmt.Printf("Adding property: %s -> %s", k, v)
		helmSetParams[k] = v
	}
	fmt.Printf("Helm Set Params %v", helmSetParams)
}

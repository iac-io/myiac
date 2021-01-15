package deploy

import (
	"fmt"
	"os"

	//"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/util"
)

type Deployer interface {
	Deploy(appName string, environment string, propertiesMap map[string]string, dryRun bool)
}

type baseDeployer struct {
	helmDeployer
	chartsPath string
}

func NewDeployerWithCharts(chartsPath string) Deployer {
	if chartsPath == "" {
		chartsPath = getBaseChartsPath()
	}
	return &baseDeployer{chartsPath: chartsPath}
}

func NewDeployer() Deployer {
	return &baseDeployer{chartsPath: getBaseChartsPath()}
}

// moneycolfrontend, moneycolserver, elasticsearch, traefik, traefik-dev, collections-api
func (bd baseDeployer) Deploy(appName string, environment string, propertiesMap map[string]string, dryRun bool) {
	helmSetParams := make(map[string]string)
	//if appName == "traefik" || appName == "traefik-dev" {
	//	helmSetParams = getNodesInternalIpsAsHelmParams()
	//}

	//TODO: Add properties to helmSetParams or values
	addPropertiesToSetParams(helmSetParams, propertiesMap)
	cmdRunner := commandline.NewEmpty()
	helmDeployer := NewHelmDeployer(bd.chartsPath, cmdRunner)
	deployment := HelmDeployment{
		AppName:       appName,
		Environment:   environment,
		HelmSetParams: helmSetParams,
		DryRun:        dryRun,
	}
	helmDeployer.Deploy(&deployment)

	// Remove once the new workflow works
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

func getBaseChartsPath() string {
	chartsPath := os.Getenv("CHARTS_PATH")
	if chartsPath != "" {
		return chartsPath
	}
	chartsPath = util.CurrentExecutableDir() + "/charts"
	return chartsPath
}

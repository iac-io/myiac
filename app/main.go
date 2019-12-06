package main

import (
	//"fmt"
	"github.com/dfernandezm/myiac/app/cli"
	//"github.com/dfernandezm/myiac/app/deploy"
	//"github.com/dfernandezm/myiac/app/docker"
	//props "github.com/dfernandezm/myiac/app/properties"
	//"github.com/dfernandezm/myiac/app/util"
)

const (
	project              = "moneycol"
	clusterZone          = "europe-west1-b"
	environment          = "dev"
	projectId            = "moneycol"
	projectRepositoryUrl = "gcr.io"
)

func main() {
	cli.BuildCli()
	//deploy.ReleaseDeployedForApp("moneycolserver")
	//cluster.InstallHelm()

	// --- Various kubernetes setups ---
	//labelNodes("elasticsearch")
	//labelNodes("applications")
}




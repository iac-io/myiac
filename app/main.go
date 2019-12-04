package main

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/cli"
	"github.com/dfernandezm/myiac/app/docker"
	props "github.com/dfernandezm/myiac/app/properties"
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
	fmt.Printf("\n\n=== MyIaC - Infrastructure as Code ===\n\n")
	cli.BuildCli()
	//cluster.InstallHelm()

	// ------ Docker workflows  -------
	//runtime := props.NewRuntime()
	//dockerWorkflows(&runtime)

	// --- Various kubernetes setups ---
	//labelNodes("elasticsearch")
	//labelNodes("applications")
}

// commit hash: git rev-parse HEAD | cut -c1-7
func dockerWorkflows(runtime *props.RuntimeProperties) {
	imageId := "b63a014d5aa6"
	commitHash := "9ff9cd6"
	appName := "moneycol-server"
	version := "0.1.0"
	dockerProps := props.DockerProperties{ProjectRepoUrl: projectRepositoryUrl, ProjectId: projectId}

	docker.TagImage(runtime, &dockerProps, imageId, commitHash, appName, version)
	docker.PushImage(runtime)
}




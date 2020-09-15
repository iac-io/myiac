package main

import (
	"github.com/dfernandezm/myiac/internal/cli"
)

const (
	project              = "moneycol"
	clusterZone          = "europe-west1-b"
	environment          = "dev"
	projectId            = "moneycol"
	projectRepositoryUrl = "eu.gcr.io"
)

func main() {
	cli.BuildCli()
}




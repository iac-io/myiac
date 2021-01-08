package main

import (
	"github.com/iac-io/myiac/internal/cli"
)

const (
	project              = "myiac"
	clusterZone          = "europe-west2-b"
	environment          = "dev"
	projectId            = "myiac"
	projectRepositoryUrl = "eu.gcr.io"
)

func main() {
	cli.BuildCli()
}

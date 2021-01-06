package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/iac-io/myiac/internal/util"
	"github.com/urfave/cli"
	"os"
)

func setupEnvironmentCmd(projectFlag *cli.StringFlag, keyPath *cli.StringFlag) cli.Command {

	providerFlag := &cli.StringFlag{
		Name: "provider",
		Usage: "Cloud Provider (gcp only for now)",
	}

	return cli.Command{
		Name:  "setupEnvironment",
		Usage: "Setup the environment with the cloud provider",
		Flags: []cli.Flag{
			providerFlag,
			projectFlag,
			keyPath,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for setupEnvironment\n")
			_ = validateBaseFlags(c)
			_ = validateStringFlagPresence("provider", c)

			providerValue := c.String("provider")
			project := c.String("project")
			keyLocation := c.String("keyPath")

			// read these values from config based on project and provider
			zone := "europe-west1-b"
			clusterName := "placeholder"
			setupProvider(providerValue, zone, clusterName, project, keyLocation)

			return nil
		},
	}
}

func setupProvider(providerValue string, zone string, clusterName string, project string, keyLocation string) {
	var provider Provider
	if providerValue == "gcp" {
		// Setup ENV Variable with the json credentials
		err := gcp.SetKeyEnvVar(util.GetGcpKeyFilePath(project))
		if err != nil {
			fmt.Printf("GCP: Could not setup Environment Variable. Error: %v", err)
			os.Exit(1)
		}
		gkeCluster := GkeCluster{zone: zone, name: clusterName}
		provider = NewGcpProvider(project, keyLocation, gkeCluster)
	} else {
		panic(fmt.Errorf("invalid provider provided: %v", providerValue))
	}

	provider.Setup()
	provider.ClusterSetup()
}
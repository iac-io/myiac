package cli

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

func setupEnvironmentCmd(providerFlag *cli.StringFlag, projectFlag *cli.StringFlag,
	keyPath *cli.StringFlag) cli.Command {

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
			zone := "europe-west2-b"
			clusterName := "placeholder"
			setupProvider(providerValue, zone, clusterName, project, keyLocation)

			return nil
		},
	}
}

func setupProvider(providerValue string, zone string, clusterName string, project string, keyLocation string) {
	var provider Provider
	if providerValue == "gcp" {
		gkeCluster := GkeCluster{zone: zone, name: clusterName}
		provider = NewGcpProvider(project, keyLocation, gkeCluster)
	} else {
		panic(fmt.Errorf("invalid provider provided: %v", providerValue))
	}

	provider.Setup()
	provider.ClusterSetup()
	log.Printf("Set local kubectl to project: %v \n", project)
}

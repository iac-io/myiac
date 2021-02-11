package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/deploy"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/urfave/cli"
	"log"
)

func installHelmCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "installHelm",
		Usage: "Install helm (tiller) in a Kubernetes cluster",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for install Helm\n")
			validateBaseFlags(c)

			project := c.String("project")
			gcp.SetupEnvironment(project)
			env := c.String("env")

			//TODO: read from project manifest
			zone := "europe-west1-b"
			gcp.SetupKubernetes(project, zone, env)

			//TODO: Need to botsrap this be full blown go solution. Disabling for now.
			//cluster.InstallHelm()

			return nil
		},
	}
}

// TODO: limit to max 7 params
func helmDeployApp(
	providerFlag *cli.StringFlag,
	projectFlag *cli.StringFlag,
	environmentFlag *cli.StringFlag,
	appName *cli.StringFlag,
	propertiesFlag *cli.StringFlag,
	keyPath *cli.StringFlag,
	zoneFlag *cli.StringFlag,
	chartsPath *cli.StringFlag,
	dryRunFlag *cli.BoolFlag) cli.Command {
	//appNameFlag := &cli.StringFlag{Name: "app, a",
	//	Usage: "The app to deploy. A helm chart with the same name must exist in the CHARTS_LOCATION"}
	//dryRunFlag := &cli.BoolFlag{Name: "dryRun", Usage: "Executes the command in dryRun mode"}
	return cli.Command{
		Name:  "deployApp",
		Usage: "Deploy an app (defined as a Helm chart from a Docker image) into a Kubernetes cluster in a given environment",
		Flags: []cli.Flag{
			providerFlag,
			projectFlag,
			environmentFlag,
			appName,
			dryRunFlag,
			keyPath,
			propertiesFlag,
			zoneFlag,
			chartsPath,
		},
		Action: func(c *cli.Context) error {
			log.Printf("Validating flags for deployApp\n")
			_ = validateStringFlagPresence("provider", c)
			_ = validateStringFlagPresence("project", c)
			_ = validateStringFlagPresence("env", c)
			_ = validateStringFlagPresence("app", c)
			//_ = validateStringFlagPresence("properties", c)
			_ = validateStringFlagPresence("keyPath", c)
			_ = validateStringFlagPresence("zone", c)
			//_ = validateStringFlagPresence("charts-path", c)

			//Set values
			provider := c.String("provider")
			project := c.String("project")
			env := c.String("env")
			dryrun := c.Bool("dry-run")
			app := c.String("app")
			properties := c.String("properties")
			key := c.String("keyPath")
			zone := c.String("zone")
			//chartsPath := c.String("charts-path")
			clusterName := project + "-" + env
			propertiesMap := readPropertiesToMap(properties)

			// Auth with provider
			if provider == "gcp" {
				//Setup ENV Variable with the json credentials
				gcp.SetKeyEnvVar(key)
			}

			// Set charts path if provided by user
			if c.String("charts-path") != "" {
				setHelmChartsPathVar(c.String("charts-path"))
			}
			// Auth with specified cluster
			cluster.SetupProvider(provider, zone, clusterName, project, key, dryrun)

			deploy.Deploy(app, env, propertiesMap, dryrun)
			return nil
		},
	}
}

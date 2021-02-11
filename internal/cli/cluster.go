package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/urfave/cli"
	"log"
)

func createClusterCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag,
	dryRunFlag *cli.BoolFlag, providerFlag *cli.StringFlag, keyPath *cli.StringFlag, tfConfigPath *cli.StringFlag, zoneFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "createCluster",
		Usage: "Create a Kubernetes cluster through Terraform",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
			dryRunFlag,
			providerFlag,
			keyPath,
			tfConfigPath,
			zoneFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for createCluster\n")
			_ = validateBaseFlags(c)
			_ = validateStringFlagPresence("provider", c)
			_ = validateStringFlagPresence("env", c)
			_ = validateStringFlagPresence("keyPath", c)
			_ = validateStringFlagPresence("zone", c)
			fmt.Printf("createCluster running with flags\n")

			project := c.String("project")
			env := c.String("env")
			dryrun := c.Bool("dry-run")
			provider := c.String("provider")
			key := c.String("keyPath")
			tfConfigPath := c.String("tfConfigPath")
			zone := c.String("zone")
			clusterName := project + "-" + env

			if provider == "gcp" {
				//Setup ENV Variable with the json credentials
				gcp.SetKeyEnvVar(key)
			}

			//TODO: pass-in variables
			err := cluster.CreateCluster(project, env, dryrun, tfConfigPath)
			if err != nil {
				log.Fatalf("Could not create cluster in project: %v. Error: %v", project, err)
			}
			cluster.SetupProvider(provider, zone, clusterName, project, key, dryrun)

			return nil

		},
	}
}

func destroyClusterCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag, providerFlag *cli.StringFlag,
	keyPath *cli.StringFlag, tfConfigPath *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "destroyCluster",
		Usage: "Destroy an existing Kubernetes cluster created previously through Terraform",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
			providerFlag,
			keyPath,
			tfConfigPath,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for destroyCluster\n")
			_ = validateBaseFlags(c)
			_ = validateStringFlagPresence("provider", c)
			_ = validateStringFlagPresence("env", c)
			_ = validateStringFlagPresence("keyPath", c)
			fmt.Printf("destroyCluster running with flags\n")

			project := c.String("project")
			env := c.String("env")
			provider := c.String("provider")
			keyPath := c.String("keyPath")
			tfConfigPath := c.String("tfConfigPath")

			if provider == "gcp" {
				//Setup ENV Variable with the json credentials
				gcp.SetKeyEnvVar(keyPath)
			}

			cluster.DestroyCluster(project, env, tfConfigPath)
			return nil
		},
	}
}

func resizeClusterCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	nodePoolsSizeFlag := &cli.StringFlag{Name: "nodePoolsSize",
		Usage: "Target size of all node pools"}
	return cli.Command{
		Name:  "resizeCluster",
		Usage: "resizeCluster to given capacity",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
			nodePoolsSizeFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for resizeCluster\n")
			_ = validateBaseFlags(c)
			_ = validateNodePoolsSize(c)

			project := c.String("project")
			env := c.String("env")
			nodePoolsSize := c.Int("nodePoolsSize")
			gcp.SetupEnvironment(project)

			//TODO: read from project manifest
			zone := "europe-west1-b"

			gcp.ResizeCluster(project, zone, env, nodePoolsSize)
			return nil
		},
	}
}

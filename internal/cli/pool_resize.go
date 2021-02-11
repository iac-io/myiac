package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/urfave/cli"
)

func resizePoolCmd(providerFlag *cli.StringFlag, projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag,
	poolNameFlag *cli.StringFlag, poolSizeFlag *cli.StringFlag, zoneFlag *cli.StringFlag,
	keyPath *cli.StringFlag, dryRunFlag *cli.BoolFlag) cli.Command {

	return cli.Command{
		Name:  "resizePool",
		Usage: "resizePool to given capacity",
		Flags: []cli.Flag{
			providerFlag,
			projectFlag,
			environmentFlag,
			poolNameFlag,
			poolSizeFlag,
			zoneFlag,
			keyPath,
			dryRunFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for resizePool\n")
			_ = validateBaseFlags(c)
			//_ = validateNodePoolsSize(c)
			_ = validateStringFlagPresence("provider", c)
			_ = validateStringFlagPresence("project", c)
			_ = validateStringFlagPresence("env", c)
			_ = validateStringFlagPresence("zone", c)
			_ = validateStringFlagPresence("keyPath", c)
			_ = validateStringFlagPresence("pool-name", c)
			_ = validateStringFlagPresence("pool-size", c)

			provider := c.String("provider")
			project := c.String("project")
			env := c.String("env")
			zone := c.String("zone")
			key := c.String("keyPath")
			poolName := c.String("pool-name")
			poolSize := c.String("pool-size")
			dryrRun := c.Bool("dry-run")

			//TODO: read from project manifest
			//log.Printf("project: %s", project)
			//log.Printf("env: %s", env)
			//log.Printf("zone: %s", zone)
			//log.Printf("resizong Node Pool: %s to %s nodes", poolName, poolSize)
			cluster.SetupProvider(provider, zone, project+"-"+env, project, key, dryrRun)
			gcp.ResizePool(project, env, poolName, poolSize, zone)
			//gcp.ResizeCluster(project, zone, env, nodePoolsSize)
			return nil
		},
	}
}

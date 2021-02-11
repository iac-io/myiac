package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/cluster"
	"github.com/urfave/cli"
)

func setupEnvironmentCmd(providerFlag *cli.StringFlag, projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag,
	keyPath *cli.StringFlag, dryRunFlag *cli.BoolFlag, zone *cli.StringFlag) cli.Command {

	return cli.Command{
		Name:  "setupEnvironment",
		Usage: "Setup the environment with the cloud provider",
		Flags: []cli.Flag{
			providerFlag,
			projectFlag,
			environmentFlag,
			keyPath,
			dryRunFlag,
			zone,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for setupEnvironment\n")
			_ = validateBaseFlags(c)
			_ = validateStringFlagPresence("provider", c)
			_ = validateStringFlagPresence("env", c)
			_ = validateStringFlagPresence("project", c)
			_ = validateStringFlagPresence("keyPath", c)
			_ = validateStringFlagPresence("zone", c)

			providerValue := c.String("provider")
			project := c.String("project")
			env := c.String("env")
			keyLocation := c.String("keyPath")
			dryrun := c.Bool("dry-run")
			zone := c.String("zone")

			clusterName := project + "-" + env
			cluster.SetupProvider(providerValue, zone, clusterName, project, keyLocation, dryrun)

			return nil
		},
	}
}

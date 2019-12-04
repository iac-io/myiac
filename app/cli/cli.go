package cli

import (
	"fmt"
	"log"
	"os"
	"github.com/dfernandezm/myiac/app/gcp"
	"github.com/dfernandezm/myiac/app/deploy"
	"github.com/urfave/cli"
)

// myiac setupEnvironment --project moneycol --env dev  [--key-location /path/to/key.json]
// myiac deploy --app traefik --project moneycol --env dev
// myiac dockerSetup --project moneycol --env dev
// myiac dockerBuild --buildPath /home/app --tag gcr.io/moneycolserver:0.1.0-abcd
// myiac dockerPush gcr.io/moneycolserver:0.1.0-abcd

func BuildCli() {
	app := cli.NewApp()
	app.Name = "myiac"
	app.Usage = "Infrastructure as code for deployments and cluster management"

	environmentFlag := &cli.StringFlag{Name: "env, e", 
										Usage: "The environment to refer to (dev,prod)"}
	projectFlag := &cli.StringFlag{Name: "project, p",
									Usage: "The project to refer to (projects folder manifests)"}

	setupEnvironment := setupEnvironmentCmd(projectFlag, environmentFlag)
	dockerSetup := dockerSetupCmd(projectFlag, environmentFlag)
	deployApp := deployAppSetup(projectFlag, environmentFlag)
	app.Commands = []*cli.Command{&setupEnvironment, &dockerSetup, &deployApp}

	err := app.Run(os.Args)
  	if err != nil {
    	log.Fatal(err)
  	}
}

func setupEnvironmentCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "setupEnvironment",
		Usage: "Setup the environment with the cloud provider (GCP is supported at the moment)",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			validateBaseFlags(c)
			fmt.Printf("setupEnvironment running with flags\n")
			gcp.SetupEnvironment()
			project := c.String("project")
			env := c.String("env")

			//TODO: read from project manifest
			zone := "europe-west1-b"
			gcp.SetupKubernetes(project, zone, env)
			return nil
		},
	}
}

func dockerSetupCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "dockerSetup",
		Usage: "Setup docker login (GCR supported at the moment)",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			validateBaseFlags(c)
			fmt.Printf("dockerSetup with flags\n")
			gcp.SetupEnvironment()
			gcp.ConfigureDocker()
			return nil
		},
	}
}

func deployAppSetup(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	appNameFlag := &cli.StringFlag{Name: "app, a", Usage: "The app to deploy. A helm chart with the same name must exist in the CHARTS_LOCATION"}

	return cli.Command{
		Name:  "deploy",
		Usage: "Setup docker login (GCR supported at the moment)",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
			appNameFlag,
		},
		Action: func(c *cli.Context) error {
			validateBaseFlags(c)
			fmt.Printf("deploy with flags flags\n")
			gcp.SetupEnvironment()
			validateStringFlagPresence("app", c)
			appToDeploy := c.String("app") 
			env := c.String("env")
			deploy.DeployApp(appToDeploy, env)
			return nil
		},
	}
}

func setupEnvironmentFromContext(c *cli.Context) {
	validateBaseFlags(c)
	gcp.SetupEnvironment()
	project := c.String("project")
	env := c.String("env")
	//TODO: read from project manifest
	zone := "europe-west1-b"
	gcp.SetupKubernetes(project, zone, env)
}

func validateBaseFlags(ctx *cli.Context) error {
	project := ctx.String("project")
	validateStringFlagPresence(project, ctx)
	
	if (project != "moneycol") {
		return cli.NewExitError("Project not supported: " + project, -1)
	}
	
	env := ctx.String("env")
	validateStringFlagPresence(env, ctx)

	if (env != "dev") {
		return cli.NewExitError("Invalid environment: " + env, -1)
	} 

	return nil
}

func validateStringFlagPresence(flagName string, ctx *cli.Context) error {
	flag := ctx.String(flagName)

	if flag == "" {
		errorMsg := fmt.Sprintf("%s parameter not provided", flag)
		return cli.NewExitError(errorMsg, -1)
	}

	return nil
}

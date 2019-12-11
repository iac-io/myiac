package cli

import (
	"fmt"
	"log"
	"os"
	"github.com/dfernandezm/myiac/app/gcp"
	"github.com/dfernandezm/myiac/app/deploy"
	"github.com/dfernandezm/myiac/app/docker"
	props "github.com/dfernandezm/myiac/app/properties"
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
	app.Usage = "Infrastructure as code for Docker/Kubernetes/Helm deployments and cluster management with Terraform (GCP support for now)"

	environmentFlag := &cli.StringFlag{Name: "env, e", 
										Usage: "The environment to refer to (dev,prod)"}
	projectFlag := &cli.StringFlag{Name: "project, p",
									Usage: "The project to refer to (projects folder manifests)"}

	setupEnvironment := setupEnvironmentCmd(projectFlag, environmentFlag)
	dockerSetup := dockerSetupCmd(projectFlag, environmentFlag)
	dockerBuild := dockerBuildCmd(projectFlag)
	deployApp := deployAppSetup(projectFlag, environmentFlag)
	app.Commands = []*cli.Command{&setupEnvironment, &dockerSetup, &deployApp, &dockerBuild}

	err := app.Run(os.Args)
  	if err != nil {
    	log.Fatal(err)
  	}
}

func setupEnvironmentCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "setupEnvironment",
		Usage: "Setup the environment with the cloud provider",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for setupEnvironment\n")
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
		Usage: "Setup docker login with a container registry (defaults to cloud provider registry)",
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
//TODO: automate GetCommitHash (git rev-parse HEAD | cut -c1-7, --git-dir /path/to/gitdir )
func dockerBuildCmd(projectFlag *cli.StringFlag) cli.Command {

	appNameFlag := &cli.StringFlag{Name: "app, a", Usage: "The container to build. Should match a repo name in registry and a Helm chart folder naming convention (moneycol-server, moneycol-frontend...)"}

	buildPathFlag := &cli.StringFlag{Name: "buildPath, bp",
									Usage: "The location of the Dockerfile"}
	commitHashFlag :=  &cli.StringFlag{Name: "commit, ch",
										Usage: "The 7 digit commit hash for the tag"}
	versionFlag :=  &cli.StringFlag{Name: "version, ch",
										Usage: "The version to be built (semver major.minor.patch)"}

	
	return cli.Command{
		Name:  "dockerBuild",
		Usage: "Build a docker image, tag it and push it to registry",
		Flags: []cli.Flag{
			projectFlag,
			buildPathFlag,
			appNameFlag,
			versionFlag,
			commitHashFlag,
		},
		Action: func(c *cli.Context) error {
			validateStringFlagPresence("project", c)
			validateStringFlagPresence("buildPath", c)
			validateStringFlagPresence("version", c)
			validateStringFlagPresence("app", c)

			fmt.Printf("dockerBuild with flags\n")
			gcp.SetupEnvironment()
			gcp.ConfigureDocker()

			project := c.String("project") 
			buildPath := c.String("buildPath")
			appName := c.String("app")
			version := c.String("version")
			commit := c.String("commit")
			
			runtime := props.NewRuntime()
			dockerProps := props.DockerProperties{ProjectRepoUrl: "gcr.io", ProjectId: project}
			docker.BuildImage(&runtime, buildPath, &dockerProps, commit, appName, version)
			docker.PushImage(&runtime)
			return nil
		},
	}
}


func deployAppSetup(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	appNameFlag := &cli.StringFlag{Name: "app, a", Usage: "The app to deploy. A helm chart with the same name must exist in the CHARTS_LOCATION"}

	return cli.Command{
		Name:  "deploy",
		Usage: "Deploy an app (defined as a Helm chart from a Docker image) into a Kubernetes cluster in a given environment",
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
	project := validateStringFlagPresence("project", ctx)
	
	if (project != "moneycol") {
		return cli.NewExitError("Project not supported: " + project, -1)
	}
	
	env := validateStringFlagPresence("env", ctx)

	if (env != "dev") {
		return cli.NewExitError("Invalid environment: " + env, -1)
	} 

	return nil
}

func validateStringFlagPresence(flagName string, ctx *cli.Context) string {
	fmt.Printf("Validating flag %s", flagName)
	flag := ctx.String(flagName)
	fmt.Printf("Read flag %s as %s", flagName, flag)

	if flag == "" {
		errorMsg := fmt.Sprintf("%s parameter not provided", flag)
		cli.NewExitError(errorMsg, -1)
		return "" 
	}

	return flag
}

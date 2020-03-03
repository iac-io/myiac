package cli

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/cluster"
	"github.com/dfernandezm/myiac/app/deploy"
	"github.com/dfernandezm/myiac/app/docker"
	"github.com/dfernandezm/myiac/app/gcp"
	props "github.com/dfernandezm/myiac/app/properties"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

const GCR_PREFIX = "eu.gcr.io"

// BuildCli myiac setupEnvironment --project moneycol --env dev  [--key-location /path/to/key.json]
// myiac deploy --app traefik --project moneycol --env dev
// myiac dockerSetup --project moneycol --env dev
// myiac dockerBuild --buildPath /home/app --tag eu.gcr.io/moneycolserver:0.1.0-abcd
// myiac dockerPush eu.gcr.io/moneycolserver:0.1.0-abcd
func BuildCli() {
	app := cli.NewApp()
	app.Name = "myiac"
	app.Usage = "Infrastructure as code for Docker/Kubernetes/Helm deployments and cluster management with Terraform (GCP support for now)"

	environmentFlag := &cli.StringFlag{Name: "env, e", Usage: "The environment to refer to (dev,prod)"}
	projectFlag := &cli.StringFlag{Name: "project, p", Usage: "The project to refer to (projects folder manifests)"}
	propertiesFlag := &cli.StringFlag{Name: "properties", Usage: "Properties for deployments"}

	setupEnvironment := setupEnvironmentCmd(projectFlag, environmentFlag)
	dockerSetup := dockerSetupCmd(projectFlag, environmentFlag)
	dockerBuild := dockerBuildCmd(projectFlag)

	createClusterCmd := createClusterCmd(projectFlag, environmentFlag)
	installHelmCmd := installHelmCmd(projectFlag, environmentFlag)
	destroyClusterCmd := destroyClusterCmd(projectFlag, environmentFlag)

	deployApp := deployAppSetup(projectFlag, environmentFlag, propertiesFlag)
	app.Commands = []*cli.Command{&setupEnvironment, &dockerSetup, &deployApp, &dockerBuild, &destroyClusterCmd, &createClusterCmd, &installHelmCmd}

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

			project := c.String("project")
			env := c.String("env")

			gcp.SetupEnvironment(project)

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
			project := c.String("project")
			gcp.SetupEnvironment(project)
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
	commitHashFlag := &cli.StringFlag{Name: "commit, ch",
		Usage: "The 7 digit commit hash for the tag"}
	versionFlag := &cli.StringFlag{Name: "version, ch",
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

			project := c.String("project")
			buildPath := c.String("buildPath")
			appName := c.String("app")
			version := c.String("version")
			commit := c.String("commit")

			gcp.SetupEnvironment(project)
			gcp.ConfigureDocker()

			runtime := props.NewRuntime()
			dockerProps := props.DockerProperties{ProjectRepoUrl: GCR_PREFIX, ProjectId: project}
			docker.BuildImage(&runtime, buildPath, &dockerProps, commit, appName, version)
			docker.PushImage(&runtime)
			return nil
		},
	}
}


func deployAppSetup(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag, propertiesFlag *cli.StringFlag) cli.Command {
	appNameFlag := &cli.StringFlag{Name: "app, a",
		Usage: "The app to deploy. A helm chart with the same name must exist in the CHARTS_LOCATION"}
	dryRunFlag := &cli.BoolFlag{Name: "dryRun", Usage: "Executes the command in dryRun mode"}
	return cli.Command{
		Name:  "deploy",
		Usage: "Deploy an app (defined as a Helm chart from a Docker image) into a Kubernetes cluster in a given environment",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
			appNameFlag,
			dryRunFlag,
			propertiesFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Deploying with flags\n")
			if err := validateBaseFlags(c); err != nil {
				fmt.Printf("Returning error %e\n",err)
				return err
			}

			project := c.String("project")
			gcp.SetupEnvironment(project)

			validateStringFlagPresence("app", c)
			appToDeploy := c.String("app")
			fmt.Printf("About to deploy %s\n", appToDeploy)
			env := c.String("env")
			fmt.Printf("Properties for deployment: %s\n", c.String("properties"))
			propertiesMap := readPropertiesToMap(c.String("properties"))
			dryRun := c.Bool("dryRun")
			deploy.Deploy(appToDeploy, env, propertiesMap, dryRun)
			return nil
		},
	}
}

func createClusterCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "createCluster",
		Usage: "Create a Kubernetes cluster through Terraform",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for createCluster\n")
			validateBaseFlags(c)
			fmt.Printf("destroyCluster running with flags\n")

			project := c.String("project")
			env := c.String("env")

			//TODO: read from project manifest
			zone := "europe-west1-b"
			gcp.SetupEnvironment(project)
	

			//TODO: pass-in variables
			cluster.CreateCluster(project, env, zone)
			return nil
		},
	}
}

func destroyClusterCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "destroyCluster",
		Usage: "Destroy an existing Kubernetes cluster created previously through Terraform",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for destroyCluster\n")
			validateBaseFlags(c)
			fmt.Printf("destroyCluster running with flags\n")

			project := c.String("project")
			gcp.SetupEnvironment(project)

			env := c.String("env")

			//TODO: read from project manifest
			zone := "europe-west1-b"
			gcp.SetupKubernetes(project, zone, env)

			cluster.DestroyCluster()
			return nil
		},
	}
}

func installHelmCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "installHelm",
		Usage: "Install helm (tiller) in a Kubernetes cluster",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for installlHelm\n")
			validateBaseFlags(c)

			project := c.String("project")
			gcp.SetupEnvironment(project)
			env := c.String("env")

			//TODO: read from project manifest
			zone := "europe-west1-b"
			gcp.SetupKubernetes(project, zone, env)

			cluster.InstallHelm()

			return nil
		},
	}
}

func setupEnvironmentFromContext(c *cli.Context) {
	validateBaseFlags(c)

	project := c.String("project")
	env := c.String("env")

	gcp.SetupEnvironment(project)

	//TODO: read from project manifest
	zone := "europe-west1-b"
	gcp.SetupKubernetes(project, zone, env)
}

func validateBaseFlags(ctx *cli.Context) error {
	project := validateStringFlagPresence("project", ctx)
	
	if project != "moneycol" {
		return cli.NewExitError("Project not supported: " + project, -1)
	}

	env := validateStringFlagPresence("env", ctx)

	if env != "dev" {
		return cli.NewExitError("Invalid environment: " + env, -1)
	} 

	return nil
}

func validateStringFlagPresence(flagName string, ctx *cli.Context) string {
	fmt.Printf("Validating flag %s\n", flagName)
	flag := ctx.String(flagName)
	fmt.Printf("Read flag %s as %s\n", flagName, flag)

	if flag == "" {
		errorMsg := fmt.Sprintf("%s parameter not provided", flag)
		err := cli.NewExitError(errorMsg, -1)
		if err != nil {
			panic(err)
		}
	}

	return flag
}

func readPropertiesToMap(properties string) map[string]string {
	propertiesMap := make(map[string]string)
	if len(properties) > 0 {
		fmt.Printf("Properties line is %s\n", properties)
		keyValues := strings.Split(properties, ",")
		for _, s := range keyValues {
			keyValue := strings.Split(s, "=")
			key := keyValue[0]
			value := keyValue[1]
			fmt.Printf("Read property: %s -> %s\n", key, value)
			propertiesMap[key] = value
		}
	}
	return propertiesMap
}

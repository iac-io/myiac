package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/deploy"
	"github.com/iac-io/myiac/internal/docker"
	"github.com/iac-io/myiac/internal/encryption"
	"github.com/iac-io/myiac/internal/gcp"
	props "github.com/iac-io/myiac/internal/properties"
	"github.com/iac-io/myiac/internal/util"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"strings"
)

const GCRPrefix = "eu.gcr.io"

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
	executionFlag := &cli.BoolFlag{Name: "noop, dry-run", Usage: "Dry Run"}

	keyPath := &cli.StringFlag{Name: "keyPath", Usage: "SA key path"}
	setupEnvironment := setupEnvironmentCmd(projectFlag, keyPath)
	dockerSetup := dockerSetupCmd(projectFlag, environmentFlag)
	dockerBuild := dockerBuildCmd(projectFlag)

	createClusterCmd := createClusterCmd(projectFlag, environmentFlag, executionFlag)
	installHelmCmd := installHelmCmd(projectFlag, environmentFlag)
	destroyClusterCmd := destroyClusterCmd(projectFlag, environmentFlag)

	deployApp := deployAppSetup(projectFlag, environmentFlag, propertiesFlag)
	resizeClusterCmd := resizeClusterCmd(projectFlag, environmentFlag)
	createSecretCmd := createSecretCmd()
	cryptCmd := cryptCmd(projectFlag)
	createCertCmd := createCertCmd()

	app.Commands = []cli.Command{
		setupEnvironment,
		dockerSetup,
		deployApp,
		dockerBuild,
		destroyClusterCmd,
		createClusterCmd,
		installHelmCmd,
		resizeClusterCmd,
		createSecretCmd,
		cryptCmd,
		createCertCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func cryptCmd(projectFlag *cli.StringFlag) cli.Command {
	modeFlag := &cli.StringFlag{
		Name: "mode, m",
		Usage: "encrypt or decrypt",
	}

	filenameWithTextFlag := &cli.StringFlag{
		Name: "filename, f",
		Usage: "Location of file with plainText to encrypt or cipherText to decrypt. " +
			"The CipherText will be written in a file with the " +
			"same name ended with .enc, the plainText file will be written with same filename ending .dec",
	}

	return cli.Command{
		Name:  "crypt",
		Usage: "Encrypt or decrypt file contents",
		Flags: []cli.Flag{
			projectFlag,
			modeFlag,
			filenameWithTextFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for crypt \n")

			_ = validateStringFlagPresence("project", c)
			_ = validateStringFlagPresence("mode", c)
			_ = validateStringFlagPresence("filename", c)

			project := c.String("project")
			mode := c.String("mode")
			filename := c.String("filename")

			gcp.SetupEnvironment(project)

			keyRingName := fmt.Sprintf("%s-keyring", project)
			keyName := fmt.Sprintf("%s-infra-key", project)
			locationId := "global"
			kmsEncrypter := gcp.NewKmsEncrypter(project, locationId, keyRingName, keyName)
			encrypter := encryption.NewEncrypter(kmsEncrypter)

			if mode != "encrypt" && mode != "decrypt" {
				return cli.NewExitError("mode can only be 'encrypt' or 'decrypt'",-1)
			}

			if mode == "encrypt" {
				encrypter.EncryptFileContents(filename)
			}

			if mode == "decrypt" {
				encrypter.DecryptFileContents(filename)
			}

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

func dockerSetupCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag) cli.Command {
	return cli.Command{
		Name:  "dockerSetup",
		Usage: "Setup docker login with a container registry (defaults to cloud provider registry)",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
		},
		Action: func(c *cli.Context) error {
			_ = validateBaseFlags(c)
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
	appNameFlag := &cli.StringFlag{
		Name: "app, a",
		Usage: "The container to build. Should match a repo name in registry " +
			"and a Helm chart folder naming convention (moneycol-server, moneycol-frontend...)",
	}
	buildPathFlag := &cli.StringFlag{Name: "buildPath, bp",
		Usage: "The location of the Dockerfile"}
	commitHashFlag := &cli.StringFlag{Name: "commit, ch", // Make sure the abbreviations don't repeat, obscure panic error happens
		Usage: "The 7 digit commit hash for the tag"}
	versionFlag := &cli.StringFlag{Name: "version, v",
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
			dockerProps := props.DockerProperties{ProjectRepoUrl: GCRPrefix, ProjectId: project}
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

func createClusterCmd(projectFlag *cli.StringFlag, environmentFlag *cli.StringFlag,
	executionFlag *cli.BoolFlag) cli.Command {
	return cli.Command{
		Name:  "createCluster",
		Usage: "Create a Kubernetes cluster through Terraform",
		Flags: []cli.Flag{
			projectFlag,
			environmentFlag,
			executionFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for createCluster\n")
			validateBaseFlags(c)
			fmt.Printf("createCluster running with flags\n")

			project := c.String("project")
			env := c.String("env")
			execflag := c.Bool("noop")

			//TODO: read from project manifest
			zone := "europe-west1-b"
			flag := &cli.StringFlag{
				Name:  project,
				Value: project,
			}
			key := &cli.StringFlag{FilePath: util.GetGcpKeyFilePath(project)}
			setupEnvironmentCmd(flag, key)

			//TODO: pass-in variables
			cluster.CreateCluster(project, env, zone, execflag)
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
			//gcp.SetupKubernetes(project, zone, env)

			cluster.DestroyCluster(project, env, zone)
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

// --- Aux functions ---

func validateBaseFlags(ctx *cli.Context) error {
	validateStringFlagPresence("project", ctx)
	return nil
}

func validateNodePoolsSize(ctx *cli.Context) error {
	fmt.Printf("Validating flag nodePoolsSize\n")
	nodePoolsSizeValue := ctx.String("nodePoolsSize")
	fmt.Printf("Flag nodePoolsSize read as %s\n", nodePoolsSizeValue)

	val, err := strconv.Atoi(nodePoolsSizeValue)

	if err != nil {
		log.Printf("Error converting %v\n", err)
		logErrorAndExit("Invalid nodePoolsSize: " + nodePoolsSizeValue)
	}

	if val >= 0 {
		// Valid, it's greater than 0
		fmt.Printf("Valid nodePoolSize %s", nodePoolsSizeValue)
		return nil
	}

	if val < 0 {
		logErrorAndExit(fmt.Sprintf("Invalid nodePoolsSize: %d", val))
	}
	
	if len(nodePoolsSizeValue) == 0 {
		logErrorAndExit("Invalid nodePoolsSize: " + nodePoolsSizeValue)
	} 
	
	return nil
}

func logErrorAndExit(errorMsg string) {
	err := cli.NewExitError(errorMsg, -1)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
}

func validateStringFlagPresence(flagName string, ctx *cli.Context) string {
	log.Printf("Validating flag %s\n", flagName)
	flag := ctx.String(flagName)
	if flag == "" {
		log.Fatalf("%s parameter not provided\n", flagName)
	}
	log.Printf("Read flag %s as %s\n", flagName, flag)
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

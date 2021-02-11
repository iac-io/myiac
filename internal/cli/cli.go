package cli

import (
	"github.com/urfave/cli"
	"log"
	"os"
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

	// multi-usage flags
	environmentFlag := &cli.StringFlag{
		Name:  "env, e",
		Usage: "The environment to refer to (dev,prod)"}
	projectFlag := &cli.StringFlag{
		Name:  "project, p",
		Usage: "The project to refer to (projects folder manifests)"}
	dryRunFlag := &cli.BoolFlag{
		Name:  "dry-run",
		Usage: "Dry Run"}
	providerFlag := &cli.StringFlag{
		Name:  "provider",
		Usage: "Select k8s provider (GCP only for now) "}
	keyPath := &cli.StringFlag{
		Name:  "keyPath",
		Usage: "SA key path"}
	tfConfigPath := &cli.StringFlag{
		Name:  "tfConfigPath",
		Usage: "Terraform Configuration Directory Path"}
	zoneFlag := &cli.StringFlag{
		Name:  "zone",
		Usage: "Cluster Zone  (example: europe-west2-b)"}
	propertiesFlag := &cli.StringFlag{
		Name:  "properties",
		Usage: "Properties for deployments"}
	poolNameFlag := &cli.StringFlag{
		Name:  "pool-name",
		Usage: "Pool Name"}
	poolSizeFlag := &cli.StringFlag{
		Name:  "pool-size",
		Usage: "New Pool Size"}
	appNameFlag := &cli.StringFlag{
		Name:  "app, a",
		Usage: "The app to deploy. A helm chart with the same name must exist in the chartsPath location"}
	chartsPath := &cli.StringFlag{
		Name:  "charts-path",
		Usage: "Helm charts Path"}

	// cmd
	setupEnvironment := setupEnvironmentCmd(
		providerFlag,
		projectFlag,
		environmentFlag,
		keyPath,
		dryRunFlag,
		zoneFlag)
	dockerSetup := dockerSetupCmd(
		projectFlag,
		environmentFlag)
	dockerBuild := dockerBuildCmd(projectFlag)
	createClusterCmd := createClusterCmd(
		projectFlag,
		environmentFlag,
		dryRunFlag,
		providerFlag,
		keyPath,
		tfConfigPath,
		zoneFlag)
	installHelmCmd := installHelmCmd(
		projectFlag,
		environmentFlag)
	destroyClusterCmd := destroyClusterCmd(
		projectFlag,
		environmentFlag,
		providerFlag,
		keyPath,
		tfConfigPath)
	deployApp := helmDeployApp(
		providerFlag,
		projectFlag,
		environmentFlag,
		appNameFlag,
		propertiesFlag,
		keyPath,
		zoneFlag,
		chartsPath,
		dryRunFlag)
	resizeClusterCmd := resizeClusterCmd(
		projectFlag,
		environmentFlag)
	resizePoolCmd := resizePoolCmd(
		providerFlag,
		projectFlag,
		environmentFlag,
		poolNameFlag,
		poolSizeFlag,
		zoneFlag,
		keyPath,
		dryRunFlag)
	createSecretCmd := createSecretCmd()
	cryptCmd := cryptCmd(projectFlag)
	createCertCmd := createCertCmd()

	// cli.Command options
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
		resizePoolCmd,
	}

	// Check if no err
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/docker"
	"github.com/iac-io/myiac/internal/gcp"
	props "github.com/iac-io/myiac/internal/properties"
	"github.com/urfave/cli"
)

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

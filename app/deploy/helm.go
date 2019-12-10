package deploy

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"github.com/dfernandezm/myiac/app/util"
)

//https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/mocking
type Release struct {
	Name       string
	Revision   int
	Updated    string
	Status     string
	Chart      string
	AppVersion string
	Namespace  string
}

type ReleasesList struct {
	Next     string
	Releases []*Release // using pointer as it becomes mutable (useful for tests)
}

// Implicit interface for commandline package, need access to those methods here
type CommandRunner interface {
	RunVoid()
	Output() string
	Setup(cmd string, args []string)
}

type helmDeployer struct {
	releases ReleasesList
	cmdRunner CommandRunner
}

func NewHelmDeployer(commandRunner CommandRunner) helmDeployer {
	return helmDeployer{ReleasesList{}, commandRunner}
}

func (hd *helmDeployer) DeployedReleasesExistsFor(appName string) bool {
	return hd.ReleaseFor(appName) != ""
}

func (hd *helmDeployer) ReleaseFor(appName string) string {
	releasesList := hd.ListReleases()
	for _, release := range releasesList.Releases {
		appNameIsPartOfChart := strings.Contains(strings.ToLower(release.Chart), appName)
		if appNameIsPartOfChart && release.Status == "DEPLOYED" {
			// It exists with the given name
			fmt.Printf("Release for app %s found. Name: %s, Status %s, Chart: %s\n",
				appName, release.Name, release.Status, release.Chart)
			return release.Name
		}
	}
	fmt.Printf("No releases found for app %s", appName)
	return ""
}

func (hd *helmDeployer) ListReleases() ReleasesList {
	cmdArgs := "list %s %s"
	argsArray := util.StringTemplateToArgsArray(cmdArgs, "--output", "json")
	hd.cmdRunner.Setup("helm", argsArray)
	hd.cmdRunner.RunVoid()
	cmdOutputJson := hd.cmdRunner.Output()
	// cmdResult := cmd.Run()
	// cmdOutputJson := cmdResult.Output()
	listReleases := hd.ParseReleasesList(cmdOutputJson)
	return listReleases
}

func (hd *helmDeployer) ParseReleasesList(jsonString string) ReleasesList {
	var listReleases ReleasesList
	jsonData := []byte(jsonString)
	err := json.Unmarshal(jsonData, &listReleases)
	if err != nil {
		log.Fatalf("Error parsing json to struct %v %v", jsonData, err)
	}
	return listReleases
}

package deploy

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"github.com/dfernandezm/myiac/app/util"
	"github.com/dfernandezm/myiac/app/commandline"
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
	Releases []Release
}

func ReleaseDeployedForApp(appName string) string {
	releasesList := ListReleases()
	for _, release := range releasesList.Releases {
		appNameIsPartOfChart := strings.Contains(strings.ToLower(release.Chart), appName)
		if appNameIsPartOfChart && release.Status == "DEPLOYED" {
			// It exists with the given name
			fmt.Printf("Release for app %s found. Name: %s, Status %s, Chart: %s",
				appName, release.Name, release.Status, release.Chart)
			return release.Name
		}
	}
	fmt.Printf("No releases found for app %s", appName)
	return ""
}

func ListReleases() ReleasesList {
	cmdArgs := "list %s %s"
	argsArray := util.StringTemplateToArgsArray(cmdArgs, "--output", "json")
	cmd := commandline.New("helm", argsArray)
	cmdResult := cmd.Run()
	cmdOutputJson := cmdResult.Output()
	listReleases := parse(cmdOutputJson)
	return listReleases
}

func parse(jsonString string) ReleasesList {
	var listReleases ReleasesList
	jsonData := []byte(jsonString)
	err := json.Unmarshal(jsonData, &listReleases)
	if err != nil {
		log.Fatalf("Error parsing json to struct %v %v", jsonData, err)
	}
	return listReleases
}

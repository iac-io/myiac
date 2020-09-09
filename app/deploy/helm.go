package deploy

import (
	"encoding/json"
	"fmt"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
	"io/ioutil"
	"log"
	"strings"
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

type helmDeployer struct {
	releases  ReleasesList
	cmdRunner commandline.CommandRunner
	chartsPath string
}

type HelmDeployment struct {
	AppName          string
	DryRun           bool
	Environment 	 string
	HelmSetParams    map[string]string // key value pairs
	HelmValuesParams []string // yaml filenames to pass as --values
}

func NewHelmDeployer(chartsPath string, commandRunner commandline.CommandRunner) *helmDeployer {
	hd := new(helmDeployer)
	hd.releases = ReleasesList{}
	hd.cmdRunner = commandRunner
	hd.chartsPath = chartsPath
	return hd
}

func (hd *helmDeployer) DeployedReleasesExistsFor(appName string) bool {
	return hd.ReleaseFor(appName) != ""
}

func (hd *helmDeployer) ReleaseFor(appName string) string {
	releasesList := hd.ListReleases()

	fmt.Println("Cleaning up FAILED releases")
	hd.DeleteFailedReleases()

	for _, release := range releasesList.Releases {
		appNameIsPartOfChart := strings.Contains(strings.ToLower(release.Chart), appName)
		if appNameIsPartOfChart && release.Status == "DEPLOYED" {
			// It exists with the given name
			fmt.Printf("Release for app %s found. " +
				"Name: %s, Status %s, Chart: %s\n",
				appName, release.Name, release.Status, release.Chart)
			return release.Name
		}
	}
	fmt.Printf("No releases found for app %s\n", appName)
	return ""
}

func (hd *helmDeployer) ListReleases() ReleasesList {
	cmdArgs := "list %s %s"
	argsArray := util.StringTemplateToArgsArray(cmdArgs, "--output", "json")
	hd.cmdRunner.Setup("helm", argsArray)
	hd.cmdRunner.RunVoid()
	cmdOutputJson := hd.cmdRunner.Output()
	listReleases := hd.ParseReleasesList(cmdOutputJson)
	return listReleases
}

func (hd *helmDeployer) DeleteFailedReleases()  {
	releasesList := hd.ListReleases()
	for _, release := range releasesList.Releases {
		if release.Status == "FAILED" {
			fmt.Printf("Deleting FAILED Helm release -> %s\n", release.Name)
			deleteHelmRelease(hd, release.Name)
		}
	}
}

func deleteHelmRelease(hd *helmDeployer, releaseName string) {
	cmdArgs := "delete %s"
	argsArray := util.StringTemplateToArgsArray(cmdArgs, releaseName)
	hd.cmdRunner.IgnoreError(true)
	hd.cmdRunner.Setup("helm", argsArray)
	hd.cmdRunner.RunVoid()
	hd.cmdRunner.IgnoreError(false)
}

func (hd *helmDeployer) ParseReleasesList(jsonString string) ReleasesList {
	var listReleases ReleasesList
	
	// If there is no releases, a single space is returned
	if jsonString == "" || len(strings.TrimSpace(jsonString)) == 0 {
		// empty releases list
		log.Printf("Empty list of releases found")
		listReleases = ReleasesList{"", []*Release{}}
	} else {
		jsonData := []byte(jsonString)
		err := json.Unmarshal(jsonData, &listReleases)
		if err != nil {
			log.Fatalf("Error parsing json to struct %v %v", jsonData, err)
		}
	}

	return listReleases
}

func (hd *helmDeployer) findChartForApp(appName string) string {
	files, err := ioutil.ReadDir(getBaseChartsPath())
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		chartFolder := file.Name()
		fmt.Printf("Checking chart folder: %s\n", chartFolder)
		appNameNormalized := normalizeName(appName)
		chartFolderNormalized := normalizeName(chartFolder)

		if appNameNormalized == chartFolderNormalized {
			fmt.Printf("Found chart for app [%s] -> [%s]\n", appNameNormalized, chartFolderNormalized)
			return getBaseChartsPath() + "/" + chartFolder
		}
	}

	log.Fatalf("Could not find chart for app %s in path %s\n", appName, getBaseChartsPath())
	return ""
}

func (hd *helmDeployer) Deploy(helmDeployment *HelmDeployment) {
	var helmArgs = ""
	var chartPathForApp = hd.findChartForApp(helmDeployment.AppName)
	var existingRelease = hd.ReleaseFor(helmDeployment.AppName)
	var action = "install"
	if existingRelease != "" {
		action = fmt.Sprintf("upgrade %s", existingRelease)
	}

	helmArgs = fmt.Sprintf("%s %s", helmArgs, action)
	helmArgs = fmt.Sprintf("%s %s", helmArgs, chartPathForApp)

	if len(helmDeployment.HelmValuesParams) > 0 {
		valuesParams := ""
		for _, filePath := range helmDeployment.HelmValuesParams {
			valuesParams += valuesParams + "--values " + filePath + " "
		}
		valuesParams = strings.TrimSpace(valuesParams)
		helmArgs = fmt.Sprintf("%s %s", helmArgs, valuesParams)
	}

	if len(helmDeployment.HelmSetParams) > 0 {
		setParams := ""
		for k, v := range helmDeployment.HelmSetParams {
			setParams += setParams + "--set " + k + "=" + v + " "
		}
		setParams = strings.TrimSpace(setParams)
		helmArgs = fmt.Sprintf("%s %s", helmArgs, setParams)
	}

	if helmDeployment.DryRun {
		helmArgs += " --debug --dry-run"
	}

	argsArray := strings.Fields(helmArgs)
	cmd := commandline.New("helm", argsArray)
	cmd.Run()
	fmt.Printf("Finished deploying app: %s\n\n", helmDeployment.AppName)
}

func normalizeName(str string) string {
	result := strings.ToLower(str)
	result = strings.TrimSpace(result)
	result = strings.Replace(result, "-", "", -1)
	return result
}
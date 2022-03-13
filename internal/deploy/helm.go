package deploy

import (
	"encoding/json"
	"fmt"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/util"
	"io/ioutil"
	"log"
	"strings"
)

//https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/mocking
type Release struct {
	Name       string
	Revision   string
	Updated    string
	Status     string
	Chart      string
	AppVersion string `json:"app_version"`
	Namespace  string
}

type file interface {
	Name() string
}


//https://talks.golang.org/2012/10things.slide#8
//https://stackoverflow.com/questions/20923938/how-would-i-mock-a-call-to-ioutil-readfile-in-go/37035375
//https://godocs.io/testing/fstest
type FileReader interface {
	ReadDir(dir string) ([]NamedFile, error)
}

type NamedFile interface {
	Name() string
}

type namedFile struct {
	name string
}

func (nf namedFile) Name() string {
	return nf.name
}

type fileReader struct {}

// ReadDir is the real implementation of the `ioutil` to read the contents
// of a directory.
//
// In this case only the name of the file is important
// so that is wrapped into a NamedFile interface (just Name() method).
// This way we make the code testable via fine-tuned interfaces
//
// See: https://stackoverflow.com/a/55133441/2128730
func (fr *fileReader) ReadDir(dir string) ([]NamedFile, error) {
	fileInfos, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	var namedFileInfos = make([]NamedFile, len(fileInfos))
	for _, fileInfo := range fileInfos {
		namedFileInfo := namedFile{
			name: fileInfo.Name(),
		}
		namedFileInfos = append(namedFileInfos, namedFileInfo)
	}
	return namedFileInfos, nil
}

type helmDeployer struct {
	cmdRunner  commandline.CommandRunner
	chartsPath string
	fileReader FileReader
}

type HelmDeployment struct {
	AppName          string
	DryRun           bool
	Environment      string
	HelmSetParams    map[string]string // key value pairs
	HelmValuesParams []string          // yaml filenames to pass as --values
}

func NewHelmDeployer(chartsPath string, commandRunner commandline.CommandRunner, outFileReader FileReader) *helmDeployer {
	hd := new(helmDeployer)
	hd.cmdRunner = commandRunner
	hd.chartsPath = chartsPath

	
	if outFileReader != nil {
		fmt.Printf("Setting up an outer fileReader %v\n", outFileReader)
		hd.fileReader = outFileReader
	} else {
		fmt.Printf("using internal fileReader \n")
		hd.fileReader = new(fileReader)
	}
	
	return hd
}

func (hd *helmDeployer) DeployedReleasesExistsFor(appName string) bool {
	return hd.ReleaseFor(appName) != ""
}

func (hd *helmDeployer) ReleaseFor(appName string) string {
	releasesList := hd.ListReleases()

	fmt.Println("Cleaning up FAILED releases")
	hd.DeleteFailedReleases()

	for _, release := range releasesList {
		appNameIsPartOfChart := strings.Contains(strings.ToLower(release.Chart), appName)
		if appNameIsPartOfChart && release.Status == "deployed" {
			// It exists with the given name
			fmt.Printf("Release for app %s found. "+
				"Name: %s, Status %s, Chart: %s\n",
				appName, release.Name, release.Status, release.Chart)
			return release.Name
		}
	}
	fmt.Printf("No releases found for app %s\n", appName)
	return ""
}

func (hd *helmDeployer) ListReleases() []*Release {
	cmdArgs := "list %s %s"
	argsArray := util.StringTemplateToArgsArray(cmdArgs, "--output", "json")
	hd.cmdRunner.Setup("helm", argsArray)
	hd.cmdRunner.RunVoid()
	cmdOutputJson := hd.cmdRunner.Output()
	listReleases := hd.ParseReleasesList(cmdOutputJson)
	return listReleases
}

func (hd *helmDeployer) DeleteFailedReleases() {
	releasesList := hd.ListReleases()
	for _, release := range releasesList {
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

func (hd *helmDeployer) ParseReleasesList(jsonString string) []*Release {
	var listReleases []*Release

	// If there is no releases, a single space is returned
	if jsonString == "" || len(strings.TrimSpace(jsonString)) == 0 {
		// empty releases list
		log.Printf("Empty list of releases found")
		listReleases = []*Release{}
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

	files, err := hd.fileReader.ReadDir(getBaseChartsPath())
	
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Printf("Base charts %v\n", getBaseChartsPath())
		fmt.Printf("Value of file %v", file)
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
	var action = "install %s"
	if existingRelease != "" {
		action = fmt.Sprintf("upgrade %s", existingRelease)
	} else {
		// in helm 3, install requires release name
		action = fmt.Sprintf(action, helmDeployment.AppName)
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
	//cmd := commandline.New("helm", argsArray)
	hd.cmdRunner.Setup("helm", argsArray)
	hd.cmdRunner.Run()
	fmt.Printf("Finished deploying app: %s\n\n", helmDeployment.AppName)
}

func normalizeName(str string) string {
	result := strings.ToLower(str)
	result = strings.TrimSpace(result)
	result = strings.Replace(result, "-", "", -1)
	return result
}

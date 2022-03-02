package deploy

import (
	"encoding/json"
	"fmt"
	"github.com/iac-io/myiac/internal/commandline"
	"testing"
)

const ExistingReleasesOutput = `[
{"name":"elastic","namespace":"default","revision":"1",
"updated":"2021-11-24 07:34:26.817149 +0000 UTC","status":"deployed",
"chart":"elasticsearch-1.0.0","app_version":"6.5.0"},
{"name":"startup-daemonset","namespace":"default","revision":"1",
"updated":"2021-11-28 18:50:57.73065 +0000 UTC","status":"deployed","chart":"startup-daemonset-1.0.0","app_version":"1.0.0"},{"name":"traefik","namespace":"default","revision":"1","updated":"2021-11-28 18:33:59.089822 +0000 UTC","status":"deployed","chart":"traefik-1.78.4","app_version":"1.7.14"}]
`

// Here we implement the CommandRunner interface with a testing mock
type mockCommandRunner struct {
	executable     string
	arguments      []string
	output         string
	suppressOutput bool
}

func (mcr *mockCommandRunner) SetSuppressOutput(suppressOutput bool) {
	mcr.suppressOutput = suppressOutput
}

func (mcr *mockCommandRunner) SetOutput(output string) {
	mcr.output = output
}

func (mcr mockCommandRunner) RunVoid() {}

func (mcr *mockCommandRunner) Output() string {
	return mcr.output
}

func (mcr mockCommandRunner) Setup(executable string, args []string) {
	mcr.executable = executable
	mcr.arguments = args
}

func (mcr mockCommandRunner) SetupWithoutOutput(executable string, args []string) {
	mcr.executable = executable
	mcr.arguments = args
}

func (mcr mockCommandRunner) IgnoreError(ignoreError bool) {}

func (mcr mockCommandRunner) Run() commandline.CommandOutput {
	fmt.Printf("current commmandline %v\n", mcr.arguments)
	return commandline.CommandOutput{Output: mcr.output}
}

func (mcr mockCommandRunner) SetupCmdLine(cmdLine string) {
	// ignored
}

// https://quii.gitbook.io/learn-go-with-tests/
// To run: go test -v
func TestReleaseDeployed(t *testing.T) {
	commandRunner := &mockCommandRunner{output: ExistingReleasesOutput}
	d := NewHelmDeployer("charts", commandRunner, nil)

	if !d.DeployedReleasesExistsFor("elastic") {
		t.Errorf("The release is deployed was incorrect, got: %v, want: %v.", false, true)
	}
}

func TestReleaseHasFailed(t *testing.T) {
	commandRunner := &mockCommandRunner{output: ""}
	d := NewHelmDeployer("charts", commandRunner, nil)

	// Given: a release (2nd one) has failed status

	releasesList := d.ParseReleasesList(ExistingReleasesOutput)
	release := releasesList[1]
	release.Status = "failed"

	existingReleasesModified, err := json.Marshal(releasesList)

	if err != nil {
		t.Errorf("Failure: error marshalling %v\n %v\n", releasesList, err)
	}

	commandRunner.SetOutput(string(existingReleasesModified))

	// When: checking if it has been deployed
	deployed := d.DeployedReleasesExistsFor("startup-daemonset")

	// Then: it shouldn't be deployed but failed
	if deployed {
		t.Errorf("The release is failed but got deployed\n")
	}
}

//TODO: should be able to use partial interface implementation for
// only the Name() part, which is what's needed

type mockFileInfo struct {
	name string
}

func (mfi mockFileInfo) Name() string {
	return mfi.name
}

type mockFileReader struct {
}

func (mfr *mockFileReader) ReadDir(dir string) ([]NamedFile, error) {
	fmt.Printf("Calling fake readdir\n")
	fileInfo := mockFileInfo{name: "app"}
	// using make([]NamedFile, 1) adds nil to the slice
	var namedFileInfos []NamedFile
	namedFileInfos = append(namedFileInfos, fileInfo)
	fmt.Printf("Slice to be returned %v\n", len(namedFileInfos))
	return namedFileInfos, nil
}

func TestHelm3RequiresNameForInstall(t *testing.T) {

	// Given a helm setup with known chart for 'app'
	commandRunner := &mockCommandRunner{output: ExistingReleasesOutput}
	fileReader := new(mockFileReader)
	d := NewHelmDeployer("charts", commandRunner, fileReader)
	

	// when deploying it
	d.Deploy(&HelmDeployment{AppName: "app", DryRun: false, Environment: "dev"})

	// it works...
	fmt.Println(commandRunner.Output())
}

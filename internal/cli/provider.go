package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/preferences"
	"github.com/iac-io/myiac/internal/util"
	"log"
	"regexp"
	"strings"
)

type Provider interface {
	Setup()
	ClusterSetup()
}

type GkeCluster struct {
	zone string
	name string
}

type ProviderFactory struct {
}

func (pf ProviderFactory) getProvider() Provider {
	prefs := preferences.DefaultConfig()
	provider := prefs.Get("provider")
	if provider == "gcp" {
		keyLocation := prefs.Get("keyLocation")
		project := prefs.Get("project")
		clusterName := prefs.Get("gke.clusterName")
		clusterZone := prefs.Get("gke.clusterZone")
		gkePrefs := GkeCluster{name: clusterName, zone: clusterZone}
		return NewGcpProvider(project, keyLocation, gkePrefs)
	} else {
		panic(fmt.Sprintf("Cloud provider not present or supported supported: %s", provider))
	}
}

type GcpProvider struct {
	projectId string
	keyLocation string
	masterSaEmail string
	gkeCluster GkeCluster
}

func NewGcpProvider(projectId string, keyLocation string, gkeCluster GkeCluster) *GcpProvider {
	validateKeyLocation(keyLocation)
	return &GcpProvider{projectId:projectId,
		keyLocation:keyLocation,
		gkeCluster:gkeCluster}
}

// Setup activate master service account from key
func (gcp *GcpProvider) Setup() {
	setupDone := gcp.checkSetup()
	if !setupDone {
		cmdLine := fmt.Sprintf("gcloud auth activate-service-account --key-file %s", gcp.keyLocation)
		cmd := commandline.NewCommandLine(cmdLine)
		cmdOutput := cmd.Run()
		gcp.masterSaEmail = extractServiceAccountEmail(cmdOutput.Output)
		gcp.savePreferences()
	}
}

func (gcp GcpProvider) checkSetup() bool {
	cmdLine := fmt.Sprintf("gcloud auth list --format json")
	cmd := commandline.NewCommandLine(cmdLine)
	cmdOutput := cmd.Run()
	authList := util.ParseArray(cmdOutput.Output)

	for _, accountAuth := range authList {
		saEmail := accountAuth["account"]
		status := accountAuth["status"]

		fmt.Printf("Checking account %s\n", saEmail)
		if status == "ACTIVE"  {
			fmt.Printf("Already authenticated for %s\n", saEmail)
			return true
		}
	}

	return false
}

func extractServiceAccountEmail(setupCmdOutput string) string {
	re := regexp.MustCompile(`\[.*\]`)
	saEmail := re.FindString(setupCmdOutput)
	saEmail = strings.Replace(saEmail, "[", "", -1)
	saEmail = strings.Replace(saEmail, "]", "", -1)
	fmt.Printf("%q\n", saEmail)
	return saEmail
}

// ClusterSetup gcloud container clusters get-credentials [cluster-name]
func (gcp GcpProvider) ClusterSetup() {
	action := "container clusters get-credentials"
	clusterName := gcp.gkeCluster.name
	zone := gcp.gkeCluster.zone
	project := gcp.projectId
	cmdLine := fmt.Sprintf("gcloud %s %s --zone %s --project %s", action, clusterName, zone, project)
	cmd := commandline.NewCommandLine(cmdLine)
	cmd.Run()
	fmt.Println("GKE setup completed")
	gcp.saveGkePreferences(clusterName, zone)
}

func (gcp GcpProvider) savePreferences() {
	prefs := preferences.DefaultConfig()
	prefs.Set("provider", "gcp")
	prefs.Set("keyLocation", gcp.keyLocation)
	prefs.Set("project", gcp.projectId)
	prefs.Set("masterSaEmail", gcp.masterSaEmail)
}

func (gcp GcpProvider) saveGkePreferences(clusterName string, zone string) {
	prefs := preferences.DefaultConfig()
	prefs.Set("gke.clusterName", clusterName)
	prefs.Set("gke.clusterZone", zone)
}

func ProviderSetup() {
	providerFactory := ProviderFactory{}
	provider := providerFactory.getProvider()
	provider.Setup()
	provider.ClusterSetup()
}

func validateKeyLocation(keyLocation string) {
	if !util.FileExists(keyLocation) {
		err := fmt.Errorf("key path is invalid %s\n", keyLocation)
		panic(err)
	}
}

func SetupProvider(providerValue string, zone string, clusterName string, project string, keyLocation string) {
	var provider Provider
	if providerValue == "gcp" {
		gkeCluster := GkeCluster{zone: zone, name: clusterName}
		provider = NewGcpProvider(project, keyLocation, gkeCluster)
	} else {
		panic(fmt.Errorf("invalid provider provided: %v", providerValue))
	}

	provider.Setup()
	provider.ClusterSetup()
	log.Printf("Set local kubectl to project: %v \n", project)
}
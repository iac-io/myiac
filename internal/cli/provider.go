package cli

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/commandline"
	"github.com/dfernandezm/myiac/internal/preferences"
	"github.com/dfernandezm/myiac/internal/util"
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
	gkeCluster GkeCluster
}

func NewGcpProvider(projectId string, keyLocation string, gkeCluster GkeCluster) *GcpProvider {
	validateKeyLocation(keyLocation)
	return &GcpProvider{projectId:projectId,
		keyLocation:keyLocation,
		gkeCluster:gkeCluster}
}

// Setup activate master service account from key
func (gcp GcpProvider) Setup() {
	cmdLine := fmt.Sprintf("gcloud auth activate-service-account --key-file %s", gcp.keyLocation)
	cmd := commandline.NewCommandLine(cmdLine)
	cmd.Run()
	gcp.savePreferences()
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
}

func (gcp GcpProvider) saveGkePreferences(clusterName string, zone string) {
	prefs := preferences.DefaultConfig()
	prefs.Set("gke.clusterName", clusterName)
	prefs.Set("gke.clusterZone", zone)
}

func validateKeyLocation(keyLocation string) {
	if !util.FileExists(keyLocation) {
		err := fmt.Errorf("key path is invalid %s\n", keyLocation)
		panic(err)
	}
}
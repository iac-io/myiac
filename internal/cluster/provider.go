package cluster

import (
	"fmt"
	"log"

	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/iac-io/myiac/internal/preferences"
)

const (
	gcpProviderName = "gcp"
)

type Provider interface {
	Setup()
	ClusterSetup()
}

func ProviderSetup() {
	providerFactory := ProviderFactory{}
	provider := providerFactory.getProvider()
	provider.Setup()
	provider.ClusterSetup()
}

func SetupProvider(providerValue string, zone string, clusterName string, project string,
	keyLocation string, dryRunFlag bool) {
	var provider Provider
	if providerValue == gcpProviderName {
		gkeCluster := GkeCluster{zone: zone, name: clusterName}
		provider = NewGcpProvider(project, keyLocation, gkeCluster)
	} else {
		panic(fmt.Errorf("invalid provider provided: %v", providerValue))
	}
	provider.Setup()
	if !dryRunFlag {
		provider.ClusterSetup()
	}

	log.Printf("Set local kubectl to project: %v \n", project)
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
	if provider == gcpProviderName {
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

type gcpProvider struct {
	projectId          string
	serviceAccountAuth gcp.Auth
	gkeCluster         GkeCluster
	prefs              preferences.Preferences
	commandRunner      commandline.CommandRunner
}

func NewGcpProvider(projectId string, keyLocation string, gkeCluster GkeCluster) *gcpProvider {
	return newGcpProvider(commandline.NewEmpty(), projectId, keyLocation, gkeCluster)
}

func newGcpProvider(commandRunner commandline.CommandRunner, projectId string, keyLocation string,
	gkeCluster GkeCluster) *gcpProvider {
	auth, err := gcp.NewServiceAccountAuth(keyLocation)

	if err != nil {
		log.Fatal(fmt.Errorf("error creating service account auth from key %s: %v", keyLocation, err))
	}

	return &gcpProvider{
		projectId:          projectId,
		serviceAccountAuth: auth,
		gkeCluster:         gkeCluster,
		prefs:              preferences.DefaultConfig(),
		commandRunner:      commandRunner,
	}
}

// Setup activates the master service account from a key file
func (gp *gcpProvider) Setup() {
	log.Printf("Setting up GCP cloud provider authentication...")
	if !gp.serviceAccountAuth.IsAuthenticated() {
		gp.serviceAccountAuth.Authenticate()
		gp.SetProject()
		gp.savePreferences()
	}
}

func (gp *gcpProvider) SetProject() {
	cmdLine := fmt.Sprintf("gcloud config set project %s", gp.projectId)
	cmd := commandline.NewCommandLine(cmdLine)
	cmd.Run()
}

// ClusterSetup gcloud container clusters get-credentials [cluster-name]
func (gp gcpProvider) ClusterSetup() {
	action := "container clusters get-credentials"
	clusterName := gp.gkeCluster.name
	zone := gp.gkeCluster.zone
	project := gp.projectId
	cmdLine := fmt.Sprintf("gcloud %s %s --zone %s --project %s", action, clusterName, zone, project)
	cmd := commandline.NewCommandLine(cmdLine)
	cmd.Run()
	fmt.Println("GKE setup completed")
	gp.saveGkePreferences(clusterName, zone)
}

func (gp gcpProvider) savePreferences() {
	prefs := gp.prefs
	prefs.Set("provider", gcpProviderName)
	prefs.Set("keyLocation", gp.serviceAccountAuth.Key().KeyFileLocation)
	prefs.Set("masterSaEmail", gp.serviceAccountAuth.Key().Email)
	prefs.Set("project", gp.projectId)
}

func (gp gcpProvider) saveGkePreferences(clusterName string, zone string) {
	prefs := gp.prefs
	prefs.Set("gke.clusterName", clusterName)
	prefs.Set("gke.clusterZone", zone)
}

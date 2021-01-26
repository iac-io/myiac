package cluster

import (
	"fmt"
	"log"

	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/iac-io/myiac/internal/preferences"
	"github.com/iac-io/myiac/internal/util"
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
	projectId  string
	saKey      *gcp.ServiceAccountKey
	gkeCluster GkeCluster
	prefs      preferences.Preferences
}

func NewGcpProvider(projectId string, keyLocation string, gkeCluster GkeCluster) *gcpProvider {
	key, err := gcp.NewServiceAccountKey(keyLocation)

	if err != nil {
		log.Fatal(fmt.Errorf("error obtaining key from location %s: %v", keyLocation, err))
	}

	return &gcpProvider{
		projectId:  projectId,
		saKey:      key,
		gkeCluster: gkeCluster,
		prefs:      preferences.DefaultConfig(),
	}
}

// Setup activate master service account from key
func (gp *gcpProvider) Setup() {
	setupDone := gp.checkSetup()
	if !setupDone {
		gp.ActivateServiceAccount()
		gp.SetProject()
		gp.savePreferences()
	}
}

func (gp *gcpProvider) SetProject() {
	cmdLine := fmt.Sprintf("gcloud config set project %s", gp.projectId)
	cmd := commandline.NewCommandLine(cmdLine)
	cmd.Run()
}

func (gp *gcpProvider) ActivateServiceAccount() {
	cmdLine := fmt.Sprintf("gcloud auth activate-service-account --key-file %s", gp.saKey.KeyFileLocation)
	cmd := commandline.NewCommandLine(cmdLine)
	cmd.Run()
}

func (gp gcpProvider) checkSetup() bool {
	providedSaEmail := gp.saKey.Email
	return isAuthenticated(providedSaEmail)
}

//TODO:  Extract

func isAuthenticated(saEmail string) bool {
	authList := listActiveAuth()
	done := isProvidedSaEmailAuthenticated(saEmail, authList)
	return done
}

func listActiveAuth() []map[string]interface{} {
	cmdLine := fmt.Sprintf("gcloud auth list --format json")
	cmd := commandline.NewCommandLine(cmdLine)
	cmdOutput := cmd.Run()
	authList := util.ParseArray(cmdOutput.Output)
	return authList
}

func isProvidedSaEmailAuthenticated(providedSaEmail string, authList []map[string]interface{}) bool {
	log.Printf("Check if already authenticated with SA: %s", providedSaEmail)
	for _, accountAuth := range authList {
		saEmail := accountAuth["account"]
		status := accountAuth["status"]

		log.Printf("Checking account %s", saEmail)

		// at this point we only allow auth using the provided service account key/email
		// if running inside GCP we would get multiple ACTIVE SAs: the ones of the service this application
		// is running on
		if status == "ACTIVE" && (saEmail == providedSaEmail) {
			log.Printf("Already authenticated for %s\n", saEmail)
			return true
		}
	}
	log.Printf("Authentication is needed for SA: %s", providedSaEmail)
	return false
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

// move to gke.go?

func (gp gcpProvider) savePreferences() {
	prefs := gp.prefs
	prefs.Set("provider", gcpProviderName)
	prefs.Set("keyLocation", gp.saKey.KeyFileLocation)
	prefs.Set("project", gp.projectId)
	prefs.Set("masterSaEmail", gp.saKey.Email)
}

func (gp gcpProvider) saveGkePreferences(clusterName string, zone string) {
	prefs := gp.prefs
	prefs.Set("gke.clusterName", clusterName)
	prefs.Set("gke.clusterZone", zone)
}

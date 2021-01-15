package cluster

import (
	"fmt"
	"log"
	"time"

	"github.com/iac-io/myiac/internal/deploy"
	"github.com/iac-io/myiac/internal/gcp"
)

const (
	clusterSinglePublicIpServiceName = "traefik-dev"
)

type GkeClusterActions interface {
	SetupNodeTerminationHandlers() error
}

type gkeClusterService struct {
	deployer   deploy.Deployer
	dnsService gcp.DNSService
	gcpProject string
	environ    string
	domainName string
}

//TODO: group domainName, gcpProject, env in one object
func NewGkeClusterService(deployer deploy.Deployer, dnsService gcp.DNSService, domainName string,
	gcpProject string, env string) GkeClusterActions {
	return &gkeClusterService{
		deployer:   deployer,
		dnsService: dnsService,
		gcpProject: gcpProject,
		environ:    env,
		domainName: domainName,
	}
}

// myiac setupEnvironment --provider gcp --project moneycol --keyPath /home/app/account.json --zone europe-west1-b --env dev
// export CF_EMAIL=
// export CF_API_KEY=
// export CHARTS_PATH=/home/app/charts
// myiac setupNodeTerminationHandlers --dnsProvider cloudflare --domain moneycol.net
func (gcs gkeClusterService) SetupNodeTerminationHandlers() error {
	// deploy service
	log.Printf("setting up service to expose single-public IP for this cluster...")
	internalIps := GetInternalIpsForNodes()
	helmSetParams := getNodesInternalIpsAsHelmParams(internalIps)
	gcs.deployer.Deploy(clusterSinglePublicIpServiceName, gcs.environ, helmSetParams, false)

	log.Printf("waiting for service to deploy...")
	time.Sleep(10 * time.Second)

	log.Printf("obtaining all public ips from cluster...")
	publicIps := GetAllPublicIps()
	aPublicIP := publicIps[0] // any public ip works for this as it's clusterIP

	// update DNS entry/es
	subdomainToUpdate := "dev." + gcs.domainName
	log.Printf("Updating dns entry %s to IP %s", subdomainToUpdate, aPublicIP)
	err := gcs.dnsService.UpsertDNSEntry(subdomainToUpdate, aPublicIP)

	if err != nil {
		return fmt.Errorf("error upserting dns entry for %s: %s", subdomainToUpdate, err)
	}

	log.Printf("DNS update completed")
	return nil
}

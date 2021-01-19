package cluster

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/iac-io/myiac/internal/deploy"
	"github.com/iac-io/myiac/internal/gcp"
)

const (
	clusterSinglePublicIpServiceName = "traefik-dev"
)

type GkeClusterActions interface {
	UpdateDnsFromClusterIps(subdomains []string) error
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

// UpdateDnsFromClusterIps this is a trick that allows having a single public IP out of a GKE cluster without using a load balancer.
// It requires traefik as the Ingress controller.

// Sequence of commands to set it up using Cloudflare as the DNS provider:
// myiac setupEnvironment --provider gcp --project moneycol --keyPath /home/app/account.json --zone europe-west1-b --env dev
// export CF_EMAIL=
// export CF_API_KEY=
// export CHARTS_PATH=/home/app/charts
// myiac updateDnsFromClusterIpsCmd --dnsProvider cloudflare --domain moneycol.net
func (gcs gkeClusterService) UpdateDnsFromClusterIps(subdomains []string) error {
	// deploy service
	log.Printf("setting up service to expose single-public IP for this cluster...")
	internalIps := GetInternalIpsForNodes()
	helmSetParams := getNodesInternalIpsAsHelmParams(internalIps)
	gcs.deployer.Deploy(clusterSinglePublicIpServiceName, gcs.environ, helmSetParams, false)

	log.Printf("waiting 5s for service to deploy...")
	time.Sleep(5 * time.Second)

	log.Printf("obtaining all public ips from cluster...")
	publicIps := GetAllPublicIps()
	aPublicIP := publicIps[0] // any public ip works for this as it's clusterIP

	// update DNS entry/es
	var subdomainsToUpdate []string
	for _, sub := range subdomains {
		if !strings.HasSuffix(sub, gcs.domainName) {
			subdomainsToUpdate = append(subdomainsToUpdate, sub+"."+gcs.domainName)
		}
	}

	log.Printf("Updating dns entries %v to IP %s", subdomainsToUpdate, aPublicIP)
	err := gcs.dnsService.UpsertDNSEntries(subdomains, aPublicIP)

	if err != nil {
		return fmt.Errorf("error upserting dns entries for %s: %s", subdomainsToUpdate, err)
	}

	log.Printf("DNS update completed")
	return nil
}

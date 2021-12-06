package cluster

import (
	"fmt"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/deploy"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/iac-io/myiac/internal/util"
	"log"
	"strings"
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

// UpdateDnsFromClusterIps Sequence of commands to set it up using Cloudflare as the DNS provider:
//
// myiac setupEnvironment --provider gcp --project moneycol --keyPath /home/app/account.json --zone europe-west1-b --env dev
// export CF_EMAIL=
// export CF_API_KEY=
// export CHARTS_PATH=/home/app/charts
// myiac updateDnsFromClusterIpsCmd --dnsProvider cloudflare --domain moneycol.net
//
func (gcs gkeClusterService) UpdateDnsFromClusterIps(subdomains []string) error {

	log.Printf("setting up service to expose single-public IP for this cluster...")
	ip := FindIngressControllerNode()

	log.Printf("Ingress Controller Node IP is %s", ip)

	// update DNS entry/es
	var subdomainsToUpdate []string
	for _, sub := range subdomains {
		if !strings.HasSuffix(sub, gcs.domainName) {
			subdomainsToUpdate = append(subdomainsToUpdate, sub+"."+gcs.domainName)
		}
	}

	log.Printf("Updating dns entries %v to IP %s", subdomainsToUpdate, ip)
	err := gcs.dnsService.UpsertDNSEntries(subdomains, ip)

	if err != nil {
		return fmt.Errorf("error upserting dns entries for %s: %s", subdomainsToUpdate, err)
	}

	log.Printf("DNS update completed")
	return nil
}

func FindIngressControllerNode() string {
	cmd11 := commandline.New("kubectl",
		[]string{"get", "pods", "-l", "app=traefik", "--template",
		"'{{range .items}}{{.metadata.name}}{{end}}'", "--field-selector", "status.phase=Running"})
	output1 := cmd11.Run()
	podName := strings.TrimSpace(output1.Output)
	podName = strings.ReplaceAll(podName, "'","")
	cmd2 := fmt.Sprintf("kubectl get pod %s --template '{{.spec.nodeName}}'", podName)
	cmd22 := commandline.NewCommandLine(cmd2)
	output2 := cmd22.Run()
	nodeName := strings.TrimSpace(output2.Output)
	nodeName = strings.ReplaceAll(nodeName, "'", "")
	cmd3 := fmt.Sprintf("kubectl get node %s -o json", nodeName)
	cmd33 := commandline.NewCommandLine(cmd3)
	output3 := cmd33.Run()
	nodeJson := util.Parse(output3.Output)
	return findNodeIp(nodeJson)
}

func findNodeIp(nodeJson map[string]interface{}) string {
	indexOfAddress := 1
	status := util.GetJsonObject(nodeJson, "status")
	addresses := util.GetJsonArray(status, "addresses")

	//TODO: find by type Internal/External
	ip := util.GetStringValue(addresses[indexOfAddress], "address")
	return ip
}

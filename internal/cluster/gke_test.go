package cluster

import (
	"fmt"
	"github.com/iac-io/myiac/internal/deploy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)


type fakeDeployer struct {
	mock.Mock
}

func (fd *fakeDeployer) Deploy(appName string, environment string, propertiesMap map[string]string, dryRun bool) {
	fd.Called(appName, environment, propertiesMap, dryRun)
}

type fakeDnsService struct {
	mock.Mock
}

func (fds *fakeDnsService) UpsertDNSEntry(dnsName string, ipAddress string) error {
	args := fds.Called(dnsName, ipAddress)
	return args.Error(0)
}

func (fds *fakeDnsService) UpsertDNSEntries(dnsEntries []string, ipAddress string) error {
	args := fds.Called(dnsEntries, ipAddress)
	return args.Error(0)
}

func TestGetAllPublicIPs(t *testing.T) {
	deployer := deploy.NewDeployer()
	dnsService := new(fakeDnsService)
	dnsService.On("UpsertDNSEntry").Return(nil)
	dnsService.On("UpsertDNSEntries").Return(nil)

	gke := NewGkeClusterService(deployer, dnsService, "moneycol.net", "moneycol", "dev")
	subdomains := []string{"collections"}
	err := gke.UpdateDnsFromClusterIps(subdomains)

	assert.Nil(t, err)
}

// go test -run TestFindIngressIp ./...
func TestFindIngressIp(t *testing.T) {
	t.Skip("This calls real kubectl commands")
	ip := FindIngressControllerNode()
	fmt.Printf("IP %s",ip)
	assert.NotNil(t, ip)
}
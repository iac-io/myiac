package cloudflare

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const Skip = true

type fakeCloudflareDNSOperations struct {
	mock.Mock
	dnsRecords map[string]cloudflare.DNSRecord
}

func (fcd *fakeCloudflareDNSOperations) ZoneIDByName(zoneName string) (string, error) {
	args := fcd.Called(zoneName)
	return args.String(0), args.Error(1)
}

func (fcd *fakeCloudflareDNSOperations) DNSRecords(zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord, error) {
	args := fcd.Called(zoneID, rr)
	return args.Get(0).([]cloudflare.DNSRecord), args.Error(1)
}

func (fcd *fakeCloudflareDNSOperations) DNSRecord(zoneID, recordID string) (cloudflare.DNSRecord, error) {
	args := fcd.Called(zoneID, recordID)
	return args.Get(0).(cloudflare.DNSRecord), args.Error(1)
}

func (fcd *fakeCloudflareDNSOperations) CreateDNSRecord(
	zoneID string,
	rr cloudflare.DNSRecord) (*cloudflare.DNSRecordResponse, error) {

	args := fcd.Called(zoneID, rr)
	if fcd.dnsRecords == nil {
		fcd.dnsRecords = make(map[string]cloudflare.DNSRecord)
	}
	fcd.dnsRecords[zoneID] = rr

	return args.Get(0).(*cloudflare.DNSRecordResponse), args.Error(1)
}

func (fcd *fakeCloudflareDNSOperations) UpdateDNSRecord(zoneID, recordID string, rr cloudflare.DNSRecord) error {
	args := fcd.Called(zoneID, recordID, rr)
	return args.Error(0)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	if !Skip {
		log.Printf("Setting up...")
		deleteZone("test.net")
	} else {
		log.Printf("Setup: SKIP")
	}
}

func shutdown() {
	if !Skip {
		log.Printf("Cleaning up...")
		deleteZone("test.net")
	} else {
		log.Printf("Shutdown: SKIP")
	}

}

func TestUnitCreateDNS(t *testing.T) {
	// given
	zoneName := "test.net"
	zoneID := "anId"
	dnsName := "testing.test.net"
	ipAddress := "1.1.1.1"
	recordID := "aRecordID"

	cfDNSMock := mockCloudflareSdk(zoneName, zoneID, dnsName, ipAddress, recordID)

	cfClient := newCfDnsClient(zoneName, cfDNSMock)

	err := cfClient.CreateDNS(dnsName, ipAddress)

	if err != nil {
		log.Fatalf("error creating DNS %v", err)
	}

	data, err := cfClient.DataForDNS(dnsName)

	assert.Nil(t, err)
	assert.Equal(t, ipAddress, data)
}

func TestBaseSetupFromApiKey(t *testing.T) {
	t.Skip("Cloudflare setup is required to run this test")
	apiKey := "fake-api-key"
	accountEmail := "account@email.com"

	api, err := cloudflare.New(apiKey, accountEmail)

	if err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, api)
}

// Learning Test: how to create a zone directly using the CF API client
// Single run: go test -v -run TestCreateZone
func TestCreateZone(t *testing.T) {
	t.Skip("Cloudflare setup is required to run this test")
	api, err := cloudflare.New(os.Getenv(cfApiKeyEnvironmentVariableName), os.Getenv(cfEmailEnvironmentVariableName))
	if err != nil {
		log.Fatal(err)
	}

	account := cloudflare.Account{}

	zone, err := api.CreateZone("test.net", false, account, "")

	if err != nil {
		log.Fatalf("Failed to create zone %v", err)
	}

	zoneDetails, err := api.ZoneDetails(zone.ID)

	if err != nil {
		log.Fatalf("Failed to get zone details %v", err)
	}

	log.Printf("Created zone %s", zoneDetails.Name)
	assert.Equal(t, "test.net", zoneDetails.Name)

	api.DeleteZone(zone.ID)
	zoneDetails, err = api.ZoneDetails(zone.ID)
	assert.NotNil(t, err)
}

func TestUpdateDNS(t *testing.T) {
	t.Skip("Cloudflare setup is required to run this test")
	zoneName := "test.net"
	dnsName := "testing"
	originalIpAddress := "1.1.1.1"

	err := setupTestDNS(zoneName, dnsName, originalIpAddress)

	if err != nil {
		log.Fatalf("Error setting up DNS test %v", err)
	}

	cfClient, _ := NewFromEnv(zoneName)

	newIpAddress := "2.2.2.2"
	err = cfClient.UpdateDNS(dnsName, newIpAddress)

	if err != nil {
		log.Fatal(err)
	}

	data, err := cfClient.DataForDNS(dnsName)

	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, newIpAddress, data)
}

func TestCreateDNS(t *testing.T) {
	t.Skip("Cloudflare setup is required to run this test")
	zoneName := "test.net"
	dnsName := "testing"
	originalIpAddress := "1.1.1.1"

	err := setupTestDNS(zoneName, dnsName, originalIpAddress)

	if err != nil {
		log.Fatalf("Error setting up DNS test %v", err)
	}

	cfClient, _ := NewFromEnv(zoneName)

	otherDns := "anotherdns"
	otherIpAddress := "3.3.3.3"

	err = cfClient.CreateDNS(otherDns, otherIpAddress)

	if err != nil {
		log.Fatalf("error creating DNS %v", err)
	}

	data, err := cfClient.DataForDNS(otherDns)

	assert.Nil(t, err)
	assert.Equal(t, otherIpAddress, data)
}

func createZone(zoneName string) (string, error) {
	api := getCfApiClient()

	account := cloudflare.Account{}
	zone, err := api.CreateZone(zoneName, false, account, "")

	if err != nil {
		return "", fmt.Errorf("error creating Zone: %v", err)
	}

	return zone.ID, nil
}

func getCfApiClient() *cloudflare.API {
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_EMAIL"))
	if err != nil {
		log.Fatalf("error getting api client: %v", err)
	}

	return api
}

func deleteZone(zoneName string) error {
	api := getCfApiClient()

	zoneId, err := api.ZoneIDByName(zoneName)

	if err != nil {
		return err
	}

	_, err = api.DeleteZone(zoneId)

	log.Printf("Zone with name %s deleted", zoneName)

	return err
}

func createDnsRecord(zoneName string, dnsName string, ipAddress string) error {
	api := getCfApiClient()

	zoneId, err := api.ZoneIDByName(zoneName)

	if err != nil {
		log.Printf("Error getting Zone ID by Name %v", zoneId)
		return err
	}

	record := cloudflare.DNSRecord{Name: dnsName, Type: "A", Content: ipAddress, TTL: 300}

	response, err := api.CreateDNSRecord(zoneId, record)

	if err != nil {
		return fmt.Errorf("error creating DNS record %v", err)
	}

	log.Printf("Successfully created DNS Record %v", response)

	return nil
}

func setupTestDNS(zoneName string, dnsName string, ipAddress string) error {

	deleteZone(zoneName)

	_, err := createZone(zoneName)

	if err != nil {
		return err
	}

	err = createDnsRecord(zoneName, dnsName, ipAddress)

	if err != nil {
		return err
	}

	return nil
}

func mockCloudflareSdk(zoneName string, zoneID string, dnsName string, ipAddress string, recordID string) *fakeCloudflareDNSOperations {
	cfDNSMock := new(fakeCloudflareDNSOperations)
	cfDNSMock.On("ZoneIDByName", zoneName).Return(zoneID, nil)

	record := cloudflare.DNSRecord{Name: dnsName, Type: "A", Content: ipAddress, TTL: 300, Proxied: true}
	createdRecord := cloudflare.DNSRecord{Name: dnsName, Type: "A", Content: ipAddress, TTL: 300, Proxied: true, ID: recordID}

	response := &cloudflare.DNSRecordResponse{
		Result:     createdRecord,
		Response:   cloudflare.Response{},
		ResultInfo: cloudflare.ResultInfo{},
	}

	cfDNSMock.On("CreateDNSRecord", zoneID, record).Return(response, nil)

	records := []cloudflare.DNSRecord{createdRecord}

	cfDNSMock.On("DNSRecords", zoneID, cloudflare.DNSRecord{}).Return(records, nil)
	cfDNSMock.On("DNSRecord", zoneID, recordID).Return(createdRecord, nil)
	return cfDNSMock
}

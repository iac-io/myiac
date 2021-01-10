package cloudflare

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
)

const (
	cfApiKeyEnvironmentVariableName = "CF_API_KEY"
	cfEmailEnvironmentVariableName = "CF_EMAIL"
)

type CfClient interface {
	UpdateDNS(dnsName string, ipAddress string) error
	CreateDNS(dnsName string, ipAddress string) error
}

type cfClient struct {
	zoneName string
	cfApi *cloudflare.API
}

func NewWithApiKey(zoneName string, apiKey string, accountEmail string) CfClient {
	api, err := cloudflare.New(apiKey, accountEmail)
	if err != nil {
		log.Fatal(err)
	}

	return &cfClient{zoneName: zoneName, cfApi: api}
}

// CF_API_KEY from gs://xxx-keys/cf-key.dec
// CF_EMAIL
func NewFromEnv(zoneName string) CfClient {
	apiKey  := os.Getenv(cfApiKeyEnvironmentVariableName)
	accountEmail := os.Getenv(cfEmailEnvironmentVariableName)
	return NewWithApiKey(zoneName, apiKey, accountEmail)
}

func (cc cfClient) UpdateDNS(dnsName string, ipAddress string) error {

	// Fetch the zone ID
	zoneName := cc.zoneName
	zoneId, err := cc.cfApi.ZoneIDByName(zoneName)
	if err != nil {
		log.Printf("error %v \n", err)
		return err
	}

	records, err := cc.cfApi.DNSRecords(zoneId, cloudflare.DNSRecord{})
	if err != nil {
		log.Printf("error: %v\n",err)
		return err
	}

	var recordId = ""
	for _, r := range records {
		fmt.Printf("%s: %s -> %s\n", r.Name, r.ID, r.Content)
		if r.Name == dnsName + "." + zoneName {
			recordId = r.ID
		}
	}

	if recordId == "" {
		log.Printf("error: record not found for dns name %s", dnsName)
		return fmt.Errorf("error: record not found for dns name %s", dnsName)
	}

	dnsRecord, _ := cc.cfApi.DNSRecord(zoneId, recordId)
	dnsRecord.Content = ipAddress
	err = cc.cfApi.UpdateDNSRecord(zoneId, recordId, dnsRecord)

	if err != nil {
		log.Printf("error updating DNS record %v", err)
		return err
	}

	return nil
}

func (cc cfClient) CreateDNS(dnsName string, ipAddress string) error {
	zoneName := cc.zoneName
	zoneId, err := cc.cfApi.ZoneIDByName(zoneName)

	if err != nil {
		log.Printf("Error getting Zone ID by Name %v", zoneId)
		return err
	}

	record := cloudflare.DNSRecord{Name: dnsName, Type: "A", Content: ipAddress, TTL: 300}
	response, err := cc.cfApi.CreateDNSRecord(zoneId, record)

	if err != nil {
		log.Printf("Error creating DNS record %v", err)
		return err
	}

	log.Printf("Successfully created DNS Record %v", response)

	return nil
}



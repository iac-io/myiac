package cloudflare

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

const (
	cfApiKeyEnvironmentVariableName = "CF_API_KEY"
	cfEmailEnvironmentVariableName  = "CF_EMAIL"
)

// Wrapper around Cloudflare Go SDK that allows managing DNS entries. More specifically,
// create, update or obtain A records (IP address) linked with a dns name. Zone should be created beforehand
// in Cloudflare. An valid API key and email are required to use this client.
type CfClient interface {
	UpdateDNS(dnsName string, ipAddress string) error
	CreateDNS(dnsName string, ipAddress string) error
	DataForDNS(dnsName string) (string, error)
	UpsertDNSEntry(dnsName string, ipAddress string) error
}

type cfClient struct {
	zoneName string
	cfApi    *cloudflare.API
}

// NewWithApiKey creates a new client for Cloudflare passing in the zone name, API key, and account email
func NewWithApiKey(zoneName string, apiKey string, accountEmail string) (CfClient, error) {
	api, err := cloudflare.New(apiKey, accountEmail)
	if err != nil {
		return nil, fmt.Errorf("error creating Cloudflare client")
	}

	return &cfClient{zoneName: zoneName, cfApi: api}, nil
}

// CF_API_KEY from gs://xxx-keys/cf-key.dec
// CF_EMAIL
func NewFromEnv(zoneName string) (CfClient, error) {
	apiKey := os.Getenv(cfApiKeyEnvironmentVariableName)
	accountEmail := os.Getenv(cfEmailEnvironmentVariableName)
	return NewWithApiKey(zoneName, apiKey, accountEmail)
}

func (cc cfClient) UpdateDNS(dnsName string, ipAddress string) error {
	zoneName := cc.zoneName

	zoneId, err := cc.cfApi.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("error getting id by name: %s", err)
	}

	records, err := cc.cfApi.DNSRecords(zoneId, cloudflare.DNSRecord{})
	if err != nil {
		return fmt.Errorf("cannot read DNS records %s", err)
	}

	var subdomain = dnsName
	if !strings.HasSuffix(dnsName, zoneName) {
		subdomain = dnsName + "." + zoneName
	}

	var recordId = ""
	for _, r := range records {
		fmt.Printf("%s: %s -> %s\n", r.Name, r.ID, r.Content)
		if r.Name == subdomain {
			recordId = r.ID
		}
	}

	if recordId == "" {
		return fmt.Errorf("error: record not found for dns name %s", dnsName)
	}

	dnsRecord, _ := cc.cfApi.DNSRecord(zoneId, recordId)
	dnsRecord.Content = ipAddress
	dnsRecord.Proxied = true

	log.Printf("About to update DNS record  %s  with %s", subdomain, ipAddress)
	err = cc.cfApi.UpdateDNSRecord(zoneId, recordId, dnsRecord)

	if err != nil {
		return fmt.Errorf("error updating DNS record %v", err)
	}

	return nil
}

func (cc cfClient) CreateDNS(dnsName string, ipAddress string) error {
	zoneName := cc.zoneName
	zoneId, err := cc.cfApi.ZoneIDByName(zoneName)

	if err != nil {
		return fmt.Errorf("error getting Zone ID by Name %s: %s", zoneId, err)
	}

	record := cloudflare.DNSRecord{Name: dnsName, Type: "A", Content: ipAddress, TTL: 300, Proxied: true}
	response, err := cc.cfApi.CreateDNSRecord(zoneId, record)

	if err != nil {
		return fmt.Errorf("error creating DNS record %s", err)
	}

	log.Printf("Successfully created DNS Record %v", response)

	return nil
}

func (cc cfClient) DataForDNS(dnsName string) (string, error) {
	zoneName := cc.zoneName
	zoneId, err := cc.cfApi.ZoneIDByName(zoneName)

	if err != nil {
		return "", err
	}

	records, err := cc.cfApi.DNSRecords(zoneId, cloudflare.DNSRecord{})
	if err != nil {
		return "", fmt.Errorf("cannot read DNS records %s", err)
	}

	var subdomain = dnsName
	if !strings.HasSuffix(dnsName, zoneName) {
		subdomain = dnsName + "." + zoneName
	}

	var recordId = ""
	log.Printf("Checking record for subdomain: %s", subdomain)
	for _, r := range records {
		fmt.Printf("%s: %s -> %s\n", r.Name, r.ID, r.Content)
		if r.Name == subdomain {
			log.Printf("Found record name %s for %s, recordId: %s", r.Name, subdomain, r.ID)
			recordId = r.ID
		}
	}

	if recordId == "" {
		log.Printf("error: record not found for dns name %s", dnsName)
		return "", fmt.Errorf("error: record not found for dns name %s", dnsName)
	}

	dnsRecord, _ := cc.cfApi.DNSRecord(zoneId, recordId)

	log.Printf("Content of DNS record %s", dnsRecord.Content)

	return dnsRecord.Content, nil
}

func (cc cfClient) UpsertDNSEntry(dnsName string, ipAddress string) error {

	dataForDns, err := cc.DataForDNS(dnsName)

	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		if strings.Contains(errStr, "not found") {
			if errDns := cc.CreateDNS(dnsName, ipAddress); errDns != nil {
				return errDns
			}
		} else {
			return err
		}
	}

	if dataForDns == "" {
		if errDns := cc.CreateDNS(dnsName, ipAddress); errDns != nil {
			return errDns
		}
	}

	// update
	if err = cc.UpdateDNS(dnsName, ipAddress); err != nil {
		return err
	}

	return nil
}

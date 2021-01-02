package cloudflare

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
)

func Example() {
	// Construct a new API object
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch user details on the account
	u, err := api.UserDetails()
	if err != nil {
		log.Fatal(err)
	}
	// Print user details
	fmt.Println(u)

	// Fetch the zone ID
	id, err := api.ZoneIDByName("moneycol.net") // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	// Fetch zone details
	zone, err := api.ZoneDetails(id)
	if err != nil {
		log.Fatal(err)
	}
	// Print zone details
	fmt.Println(zone)

	records, err := api.DNSRecords(id, cloudflare.DNSRecord{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, r := range records {
		fmt.Printf("%s: %s -> %s\n", r.Name, r.ID, r.Content)
	}

	fmt.Printf("Updating IP address...")

	updateDnsIp(api, id, "collections.moneycol.net", "35.241.212.166")

	fmt.Printf("After update")

	records2, err2 := api.DNSRecords(id, cloudflare.DNSRecord{})
	if err2 != nil {
		fmt.Println(err)
		return
	}

	for _, r := range records2 {
		fmt.Printf("%s: %s -> %s\n", r.Name, r.ID, r.Content)
	}
}

func updateDnsIp(cfApi *cloudflare.API, zoneId string, dnsName string, ipAddress string) {
	records, err := cfApi.DNSRecords(zoneId, cloudflare.DNSRecord{})
	if err != nil {
		fmt.Println(err)
		return
	}

	var recordId string = ""
	for _, r := range records {
		fmt.Printf("%s: %s -> %s\n", r.Name, r.ID, r.Content)
		if r.Name == dnsName {
			recordId = r.ID
		}
	}
	if recordId == ""{
		log.Fatalf("Not found record for dns name %s", dnsName)
		return
	}
	dnsRecord, _ := cfApi.DNSRecord(zoneId, recordId)
	dnsRecord.Content = ipAddress
	err2 := cfApi.UpdateDNSRecord(zoneId, recordId, dnsRecord)

	if err2 != nil {
		log.Fatalf("error updating DNS record", err)
	}
}



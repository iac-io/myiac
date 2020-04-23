package gcp

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/dns/v1"
)

// GoogleCloudDNSService is the service that allows to create or update dns records
type GoogleCloudDNSService struct {
	service *dns.Service
	project string
	zone    string
}

// NewGoogleCloudDNSService Creates a GCP cloud dns service
// zone is the managedZone name (should've been created beforehand)
func NewGoogleCloudDNSService(project, zone string) *GoogleCloudDNSService {

	fmt.Printf("Creating new GoogleCloudDNSService for project %v and zone %v\n", project, zone)

	ctx := context.Background()

	googleClient, err := google.DefaultClient(ctx, dns.NdevClouddnsReadwriteScope)

	if err != nil {
		log.Fatal().Err(err).Msg("Creating google cloud client failed")
	}

	dnsService, err := dns.New(googleClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Creating google cloud dns service failed")
	}

	return &GoogleCloudDNSService{
		service: dnsService,
		project: project,
		zone:    zone,
	}
}

// GetDNSRecordByName returns the record sets matching name and type
// dnsRecordType is A, CNAME...
// dnsRecordName is the DNS itself (dev.moneycol.ml, ...)
func (dnsService *GoogleCloudDNSService) GetDNSRecordByName(dnsRecordType, dnsRecordName string) (records []*dns.ResourceRecordSet) {

	records = make([]*dns.ResourceRecordSet, 0)

	req := dnsService.service.ResourceRecordSets.List(dnsService.project, dnsService.zone).Name(fmt.Sprintf("%v.", dnsRecordName)).Type(dnsRecordType)

	err := req.Pages(context.Background(), func(page *dns.ResourceRecordSetsListResponse) error {
		records = page.Rrsets
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msgf("Failed retrieving records")
	}

	return
}

// UpsertDNSRecord either updates or creates a dns record.
// dnsRecordContent is usually and IP address or an alias to another service (A and CNAME records)
func (dnsService *GoogleCloudDNSService) UpsertDNSRecord(dnsRecordType, dnsRecordName, dnsRecordContent string) (err error) {

	// retrieve records in case they exist
	records := dnsService.GetDNSRecordByName(dnsRecordType, dnsRecordName)

	change := dns.Change{
		Additions: []*dns.ResourceRecordSet{
			{
				Name: fmt.Sprintf("%v.", dnsRecordName),
				Type: dnsRecordType,
				Ttl:  20,
				Rrdatas: []string{
					dnsRecordContent,
				},
				SignatureRrdatas: []string{},
				Kind:             "dns#resourceRecordSet",
			},
		},
	}

	if len(records) > 0 {
		// updating a record is done by deleting the current ones and adding the new one
		change.Deletions = records
	}

	resp, err := dnsService.service.Changes.Create(dnsService.project, dnsService.zone, &change).Context(context.Background()).Do()

	if err != nil {
		log.Fatal().Interface("Error %v", err)
		return err
	}

	log.Info().Interface("response", resp).Msgf("Response from google cloud dns api")

	return
}
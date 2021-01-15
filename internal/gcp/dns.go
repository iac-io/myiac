package gcp

import (
	"fmt"
	"strings"

	"github.com/iac-io/myiac/services"

	"github.com/iac-io/myiac/internal/cloudflare"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
)

const (
	dnsProviderGcp        = "gcp"
	dnsProviderCloudflare = "cloudflare"
)

type DNSService interface {
	UpsertDNSEntry(dnsName string, ipAddress string) error
	UpsertDNSEntries(dnsEntries []string, ipAddress string) error
}

type DNSChangeRequest interface {
	DNSProvider() string
	DomainName() string
}

type dnsChangeRequest struct {
	dnsProvider string
	domainName  string
}

func NewDNSChangeRequest(dnsProvider string, domainName string, serviceName string) (DNSChangeRequest, error) {
	if (dnsProvider != dnsProviderGcp) && (dnsProvider != dnsProviderCloudflare) {
		return nil, fmt.Errorf("error: unknown dns provider %s", dnsProvider)
	}
	serviceProps := services.NewServiceProps(serviceName)
	if !strings.HasSuffix(domainName, serviceProps.ClusterDomainName) {
		return nil, fmt.Errorf("error: domain name not supported %s", domainName)
	}

	return &dnsChangeRequest{domainName: domainName, dnsProvider: dnsProvider}, nil
}

func (dcr dnsChangeRequest) DomainName() string {
	return dcr.domainName
}

func (dcr dnsChangeRequest) DNSProvider() string {
	return dcr.dnsProvider
}

// GoogleCloudDNSService is the service that allows to create or update dns records
type GoogleCloudDNSService struct {
	service *dns.Service
	project string
	zone    string
}

func NewDNSService(dnsProvider string, serviceName string) (DNSService, error) {
	props := services.NewServiceProps(serviceName)
	if dnsProvider == dnsProviderGcp {
		log.Printf("Setting up GCP DNS provider")
		return NewGoogleCloudDNSService(props.GcpProjectId, props.GkeClusterZone), nil
	} else if dnsProvider == dnsProviderCloudflare {
		log.Printf("Setting up cloudflare DNS provider %s", props.CloudflareZone)
		service, err := cloudflare.NewFromEnv(props.CloudflareZone)
		if err != nil {
			return nil, err
		}
		return service, nil
	} else {
		return nil, fmt.Errorf("error: DNS provider not recognised %s", dnsProvider)
	}
}

// NewGoogleCloudDNSService Creates a GCP cloud dns service
// zone is the managedZone name (should've been created beforehand)
func NewGoogleCloudDNSService(project, zone string) *GoogleCloudDNSService {

	fmt.Printf("Creating new GoogleCloudDNSService for project %v and zone %v", project, zone)

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
			&dns.ResourceRecordSet{
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

	return nil
}

func (dnsService *GoogleCloudDNSService) UpsertDNSEntry(dnsName string, ipAddress string) error {
	return dnsService.UpsertDNSRecord("A", dnsName, ipAddress)
}

func (dnsService *GoogleCloudDNSService) UpsertDNSEntries(dnsEntries []string, ipAddress string) error {

	log.Info().Msg("About to update/insert DNS entries to IP " + ipAddress)

	var errors []error
	for _, dnsEntry := range dnsEntries {
		err := dnsService.UpsertDNSRecord("A", dnsEntry, ipAddress)
		if err != nil {
			fmt.Printf("Error updating DNS entry %s: %v", dnsEntry, err)
			errors = append(errors, err)
		} else {
			log.Info().Msg("DNS entry " + dnsEntry + " updated")
		}
	}

	if len(errors) > 0 {
		fmt.Printf("errors ocurred during upsert of DNS entries %v", errors)
		return errors[0]
	}

	return nil
}

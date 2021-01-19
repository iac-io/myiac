package services

const serviceMoneycol = "moneycol"

type ServiceProps struct {
	GcpProjectId      string
	GkeClusterZone    string
	CloudflareZone    string
	ClusterDomainName string
	DnsEntries        []string
}

var serviceProps *ServiceProps

func NewServiceProps(serviceName string) *ServiceProps {
	if serviceName == serviceMoneycol {
		if serviceProps == nil {
			serviceProps = &ServiceProps{
				GcpProjectId:      "moneycol",
				GkeClusterZone:    "europe-west1-b",
				CloudflareZone:    "moneycol.net",
				ClusterDomainName: "moneycol.net",
				DnsEntries:        []string{"dev", "collections", "graphql-dev"},
			}
		}
		return serviceProps
	} else {
		return nil
	}
}

package services

type ServiceProps struct {
	GcpProjectId      string
	GkeClusterZone    string
	CloudflareZone    string
	ClusterDomainName string
}

var serviceProps *ServiceProps

func NewServiceProps(serviceName string) *ServiceProps {
	if serviceName == "moneycol" {
		if serviceProps == nil {
			serviceProps = &ServiceProps{
				GcpProjectId:      "moneycol",
				GkeClusterZone:    "europe-west1-b",
				CloudflareZone:    "moneycol.net",
				ClusterDomainName: "moneycol.net",
			}
		}
		return serviceProps
	} else {
		return nil
	}
}

package cli

import (
	"fmt"

	"github.com/iac-io/myiac/services"

	"github.com/iac-io/myiac/internal/cluster"
	"github.com/iac-io/myiac/internal/deploy"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/urfave/cli"
)

// updateDnsFromClusterIpsCmd
//
// myiac setupEnvironment --provider gcp --project moneycol --keyPath /home/app/account.json --zone europe-west1-b --env dev
// DEPRECATED
func updateDnsFromClusterIpsCmdOlder() cli.Command {

	dnsProvider := &cli.StringFlag{
		Name:  "dnsProvider, dp",
		Usage: "DNS provider to update ('gcp', 'cloudflare')",
	}

	domainName := &cli.StringFlag{
		Name:  "domain, d",
		Usage: "Domain name to update (no subdomain, i.e 'moneycol.net')",
	}

	return cli.Command{
		Name:  "updateDnsWithClusterIps",
		Usage: "Setup to update dns on node termination/preemption",
		Flags: []cli.Flag{
			dnsProvider,
			domainName,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for nodeTerminationHandlers \n")

			validatePreemptionSetupFlags(c)
			dnsProvider := c.String("dnsProvider")
			domainName := c.String("domain")

			cluster.ProviderSetup()

			_, err := gcp.NewDNSChangeRequest(dnsProvider, domainName, "moneycol")

			if err != nil {
				return fmt.Errorf("error: invalid DNS change request %s", err)
			}

			projectId := "moneycol"
			env := "dev"

			deployer := deploy.NewDeployer()
			dnsService, _ := gcp.NewDNSService(dnsProvider, "moneycol")
			gkeService := cluster.NewGkeClusterService(deployer, dnsService, domainName, projectId, env)

			dnsEntries := services.NewServiceProps("moneycol").DnsEntries
			err = gkeService.UpdateDnsFromClusterIps(dnsEntries)

			if err != nil {
				return fmt.Errorf("error setting up node termination handlers: %s", err)
			}

			return nil
		},
	}
}

func validatePreemptionSetupFlags(c *cli.Context) {
	_ = validateStringFlagPresence("dnsProvider", c)
	_ = validateStringFlagPresence("domain", c)
}

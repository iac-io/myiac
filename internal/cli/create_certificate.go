package cli

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/secret"
	"github.com/dfernandezm/myiac/internal/ssl"
	"github.com/urfave/cli"
	"log"
)

func createCertCmd() cli.Command {

	keyPathFlag := &cli.StringFlag{
		Name: "keyPath, k",
		Usage: "Location of file with private key",
	}

	certPathFlag := &cli.StringFlag{
		Name: "certPath, c",
		Usage: "Cert path flag",
	}

	domainNameFlag := &cli.StringFlag{
		Name: "domain, d",
		Usage: "Domain name",
	}

	return cli.Command{
		Name:  "createCert",
		Usage: "Create certificate from files",
		Flags: []cli.Flag{
			domainNameFlag,
			keyPathFlag,
			certPathFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for createCert \n")

			validateFlags(c)
			keyPath := c.String("keyPath")
			certPath := c.String("certPath")
			domainName := c.String("domain")

			log.Printf("Creating certificate for %s from %s - %s \n", domainName, certPath, keyPath)

			ProviderSetup()

			secretManager := secret.CreateKubernetesSecretManager("default")
			certificate := ssl.NewCertificate(domainName, certPath, keyPath)
			certStore := ssl.NewSecretCertStore(secretManager)
			certStore.Register(certificate)

			return nil
		},
	}
}

func validateFlags(c *cli.Context) {
	_ = validateStringFlagPresence("keyPath", c)
	_ = validateStringFlagPresence("certPath", c)
	_ = validateStringFlagPresence("domain", c)
}



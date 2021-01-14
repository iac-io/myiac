package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/iac-io/myiac/internal/gcp"
	"github.com/iac-io/myiac/internal/secret"
	"github.com/urfave/cli"
)

func createSecretCmd() cli.Command {
	secretNameFlag := &cli.StringFlag{Name: "secretName", Usage: "The name of the secret to be created in K8s"}
	saEmailFlag := &cli.StringFlag{Name: "saEmail", Usage: "The service account email whose key will be associated to the secret"}
	recreateKeyFlag := &cli.BoolFlag{Name: "recreateSaKey", Usage: "Whether or not it should recreate the SA key"}
	literalStringFlag := &cli.StringFlag{Name: "literal", Usage: "String to encode as secret, in plain text"}

	return cli.Command{
		Name:  "createSecret",
		Usage: "Create Kubernetes secret from a file",
		Flags: []cli.Flag{
			secretNameFlag,
			saEmailFlag,
			recreateKeyFlag,
			literalStringFlag,
		},
		Action: func(c *cli.Context) error {

			fmt.Printf("Create secret command with flags\n")

			// For file-based secrets
			secretName := c.String("secretName")
			saEmail := c.String("saEmail")
			recreateKey := c.Bool("recreateSaKey")

			// for literal secrets
			literal := c.String("literal")

			cluster.ProviderSetup()

			if len(saEmail) > 0 {
				createSecretForServiceAccount(saEmail, secretName, recreateKey)
			} else if len(literal) > 0 {
				fmt.Printf("Creating secret for literal string\n")
				err := createLiteralSecret(secretName, literal)
				if err != nil {
					return fmt.Errorf("error creating literal secret %v", err)
				}
			} else {
				return fmt.Errorf("no supported secret type detected")
			}

			return nil
		},
	}
}

func createSecretForServiceAccount(saEmail string, secretName string, recreateKey bool) error {
	log.Printf("Creating secret for service account %s\n", saEmail)
	saClient := gcp.NewDefaultServiceAccountClient()
	keyFilePath := fmt.Sprintf("/tmp/%s", secretName)
	err := saClient.KeyFileForServiceAccount(saEmail, recreateKey, keyFilePath)

	if err != nil {
		return fmt.Errorf("error creating key for email %s", err)
	}

	kubeSecretManager := secret.CreateKubernetesSecretManager("default")
	kubeSecretManager.CreateFileSecret(secretName, keyFilePath)

	return nil
}

func createLiteralSecret(secretName string, literal string) error {
	literalArr := strings.Split(literal, "=")
	if len(literalArr) >= 2 {
		//TODO: support multiple literals comma separated
		literalMap := make(map[string]string)
		literalMap[literalArr[0]] = literalArr[1]
		kubeSecretManager := secret.CreateKubernetesSecretManager("default")
		kubeSecretManager.CreateLiteralSecret(secretName, literalMap)
		return nil
	} else {
		return fmt.Errorf("error, literal should have key=value pairs")
	}
}

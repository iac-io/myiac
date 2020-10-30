package cli

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/cluster"
	"github.com/dfernandezm/myiac/internal/gcp"
	"github.com/dfernandezm/myiac/internal/secret"
	"github.com/urfave/cli"
	"strings"
)

func createSecretCmd() cli.Command {
	secretNameFlag := &cli.StringFlag{Name: "secretName", Usage: "The name of the secret to be created in K8s"}
	saEmailFlag := &cli.StringFlag{Name: "saEmail", Usage: "The service account email whose key will be associated to the secret"}
	recreateKeyFlag :=  &cli.BoolFlag{Name: "recreateSaKey", Usage: "Whether or not it should recreate the SA key"}
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

			ProviderSetup()

			if len(saEmail) > 0 {
				createSecretForServiceAccount(saEmail, secretName, recreateKey)
			} else if len(literal) > 0 {
				fmt.Printf("Creating secret for literal string\n")
				literalArr := strings.Split(literal, "=")

				if len(literalArr) >= 2 {
					//TODO: support multiple literals comma separated
					//TODO: add this to the kubeSecretManager
					literalMap := make(map[string]string)
					literalMap[literalArr[0]] = literalArr[1]
					cluster.CreateSecretFromLiteral(secretName, "default", literalMap)
				} else {
					return fmt.Errorf("error, literal should have key=value pairs")
				}
			} else {
				return fmt.Errorf("no supported secret type detected")
			}

			return nil
		},
	}
}

func createSecretForServiceAccount(saEmail string, secretName string, recreateKey bool) {
	fmt.Printf("Creating secret for service account %s\n", saEmail)

	saClient := gcp.NewDefaultServiceAccountClient()
	keyFilePath := fmt.Sprintf("/tmp/%s", secretName)
	err := saClient.KeyFileForServiceAccount(saEmail, recreateKey, keyFilePath)

	if err != nil {
		panic(err)
	}

	kubeSecretManager := secret.CreateKubernetesSecretManager("default")
	kubeSecretManager.CreateFileSecret(secretName, keyFilePath)
}
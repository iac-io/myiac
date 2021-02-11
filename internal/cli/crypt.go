package cli

import (
	"fmt"
	"github.com/iac-io/myiac/internal/encryption"
	"github.com/iac-io/myiac/internal/gcp"
	"github.com/urfave/cli"
)

func cryptCmd(projectFlag *cli.StringFlag) cli.Command {
	modeFlag := &cli.StringFlag{
		Name:  "mode, m",
		Usage: "encrypt or decrypt",
	}

	filenameWithTextFlag := &cli.StringFlag{
		Name: "filename, f",
		Usage: "Location of file with plainText to encrypt or cipherText to decrypt. " +
			"The CipherText will be written in a file with the " +
			"same name ended with .enc, the plainText file will be written with same filename ending .dec",
	}

	return cli.Command{
		Name:  "crypt",
		Usage: "Encrypt or decrypt file contents",
		Flags: []cli.Flag{
			projectFlag,
			modeFlag,
			filenameWithTextFlag,
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("Validating flags for crypt \n")

			_ = validateStringFlagPresence("project", c)
			_ = validateStringFlagPresence("mode", c)
			_ = validateStringFlagPresence("filename", c)

			project := c.String("project")
			mode := c.String("mode")
			filename := c.String("filename")

			gcp.SetupEnvironment(project)

			keyRingName := fmt.Sprintf("%s-keyring", project)
			keyName := fmt.Sprintf("%s-infra-key", project)
			locationId := "global"
			kmsEncrypter := gcp.NewKmsEncrypter(project, locationId, keyRingName, keyName)
			encrypter := encryption.NewEncrypter(kmsEncrypter)

			if mode != "encrypt" && mode != "decrypt" {
				return cli.NewExitError("mode can only be 'encrypt' or 'decrypt'", -1)
			}

			if mode == "encrypt" {
				encrypter.EncryptFileContents(filename)
			}

			if mode == "decrypt" {
				encrypter.DecryptFileContents(filename)
			}

			return nil
		},
	}
}

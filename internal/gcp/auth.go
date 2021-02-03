package gcp

import (
	"fmt"
	"log"

	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/util"
)

const (
	saStatusActive        = "ACTIVE"
	statusAuthField       = "status"
	emailAccountAuthField = "account"
	listAuthCmd           = "gcloud auth list --format json -q"
)

type Auth interface {
	Authenticate()
	IsAuthenticated() bool
	Key() *ServiceAccountKey
}

type ServiceAccountAuth struct {
	ServiceAccountKey *ServiceAccountKey
	commandRunner     commandline.CommandRunner
}

func NewServiceAccountAuth(keyLocation string) (Auth, error) {
	saKey, err := NewServiceAccountKey(keyLocation)

	if err != nil {
		return nil, err
	}

	return newServiceAccountAuthWithRunner(commandline.NewEmpty(), saKey)
}

// private constructor that receives commandline (useful for testing)
func newServiceAccountAuthWithRunner(commandRunner commandline.CommandRunner, accountKey *ServiceAccountKey) (Auth, error) {
	return &ServiceAccountAuth{ServiceAccountKey: accountKey, commandRunner: commandRunner}, nil
}

func (saa *ServiceAccountAuth) Authenticate() {
	keyLocation := saa.ServiceAccountKey.KeyFileLocation
	cmdLine := fmt.Sprintf("gcloud auth activate-service-account --key-file %s", keyLocation)
	saa.commandRunner.SetupCmdLine(cmdLine)
	saa.commandRunner.Run()
}

func (saa ServiceAccountAuth) IsAuthenticated() bool {
	authList := saa.listActiveAuth()
	done := saa.isServiceAccountEmailAuthenticated(authList)
	return done
}

func (saa ServiceAccountAuth) Key() *ServiceAccountKey {
	return saa.ServiceAccountKey
}

func (saa ServiceAccountAuth) listActiveAuth() []map[string]interface{} {
	cmdLine := fmt.Sprintf(listAuthCmd)
	saa.commandRunner.SetupCmdLine(cmdLine)
	cmdOutput := saa.commandRunner.Run()
	authList := util.ParseArray(cmdOutput.Output)
	return authList
}

func (saa ServiceAccountAuth) isServiceAccountEmailAuthenticated(authList []map[string]interface{}) bool {
	providedSaEmail := saa.ServiceAccountKey.Email

	log.Printf("Check if already authenticated with SA: %s", providedSaEmail)
	for _, accountAuth := range authList {
		saEmail := accountAuth[emailAccountAuthField]
		status := accountAuth[statusAuthField]

		log.Printf("Checking account %s", saEmail)

		// at this point it's only allowed / considered authentication using the provided service account key.
		// if running inside GCP there will be multiple ACTIVE SAs: the ones of the service this application
		// is running on (GKE)
		if status == saStatusActive && (saEmail == providedSaEmail) {
			log.Printf("Already authenticated for %s", saEmail)
			return true
		}
	}
	log.Printf("Authentication is needed for SA: %s", providedSaEmail)
	return false
}

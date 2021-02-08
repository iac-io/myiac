package gcp

import (
	"fmt"
	"log"

	"github.com/iac-io/myiac/internal/util"
)

type ServiceAccountKey struct {
	KeyFileLocation string
	Email           string
}

func NewServiceAccountKey(keyLocation string) (*ServiceAccountKey, error) {
	ok, err := isValidKeyLocation(keyLocation)

	if !ok {
		log.Printf("error validating key -- this is required %s", err)
		return nil, err
	}

	saEmail, err := extractSaEmailFromKey(keyLocation)

	if err != nil {
		log.Printf("error: email could not be obtained from key at location %s", saEmail)
	}

	return &ServiceAccountKey{KeyFileLocation: keyLocation, Email: saEmail}, nil
}

func extractSaEmailFromKey(keyLocation string) (string, error) {
	json, err := util.ReadFileToString(keyLocation)

	if err != nil {
		return "", fmt.Errorf("error reading key location %s", err)
	}

	keyJson := util.Parse(json)
	saEmail := util.GetStringValue(keyJson, "client_email")

	log.Printf("Service account email in JSON key is: %s", saEmail)
	return saEmail, nil
}

func isValidKeyLocation(keyLocation string) (bool, error) {
	if !util.FileExists(keyLocation) {
		return false, fmt.Errorf("key path is invalid %s\n", keyLocation)
	}
	return true, nil
}

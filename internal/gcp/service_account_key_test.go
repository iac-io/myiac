package gcp

import (
	"fmt"
	"log"
	"testing"

	"github.com/iac-io/myiac/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestValidKeyWithValidEmail(t *testing.T) {

	accountEmail := "testAccount@testProject.iam.gserviceaccount.com"
	accountKey, err := NewServiceAccountKey(testAccountKeyPath)

	if err != nil {
		log.Fatalf("error creating service account key account %s", err)
	}

	assert.Equal(t, accountEmail, accountKey.Email)
}

func TestInvalidValidKeyLocationFails(t *testing.T) {
	testAccountKeyPath := "test_account.json"
	_, err := NewServiceAccountKey(testAccountKeyPath)

	if err != nil {
		errMessage := fmt.Sprintf("%s", err)
		assert.Contains(t, errMessage, "key path is invalid")
		return
	}

	assert.Fail(t, "key location is invalid -- test should've failed")
}

// util
func WriteSaKeyToString(testAccountKeyPath string) {
	err := util.WriteStringToFile(testAccountKeyJson, testAccountKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	return
}

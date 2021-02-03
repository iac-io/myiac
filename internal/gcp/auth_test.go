package gcp

import (
	"log"
	"testing"

	"github.com/iac-io/myiac/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	statusActiveJson = `
	[
	  {
		"account": "testAccount@testProject.iam.gserviceaccount.com",
		"status": "ACTIVE"
	  },
	  {
		"account": "otherAccount@testProject.iam.gserviceaccount.com",
		"status": ""
	  }
	]
	`

	statusNotActiveJson = `
	[
	  {
		"account": "testAccount@testProject.iam.gserviceaccount.com",
		"status": ""
	  },
	  {
		"account": "otherAccount@testProject.iam.gserviceaccount.com",
		"status": "ACTIVE"
	  }
	]
	`
	statusEmptyJson = "[]"

	//TODO: throws unmarshalling error, should be captured?
	statusEmptyString = ""
)

//func TestAuthenticatesWhenValidSa(t *testing.T) {
//	testAccountKeyPath := "../../testdata/test_account.json"
//	accountKey, _ := NewServiceAccountKey(testAccountKeyPath)
//
//	cmdLine := testutil.FakeCommandRunner("test-output")
//	saAuth, err := newServiceAccountAuthWithRunner(cmdLine, accountKey)
//
//	if err != nil {
//		log.Fatalf("error creating auth -- %s", err)
//	}
//
//	saAuth.Authenticate()
//
//	cmdLineRun := cmdLine.GetCmdLines()[0]
//	expectedCmdLine := fmt.Sprintf("gcloud auth activate-service-account --key-file %s", testAccountKeyPath)
//	assert.Equal(t, expectedCmdLine, cmdLineRun)
//}

func TestChecksAuthDone(t *testing.T) {
	testAccountKeyPath := "../../testdata/test_account.json"
	accountKey, _ := NewServiceAccountKey(testAccountKeyPath)

	cmdLine := testutil.FakeCommandRunner("default-output")
	cmdLine.FakeCommand(listAuthCmd, statusActiveJson)

	saAuth, err := newServiceAccountAuthWithRunner(cmdLine, accountKey)

	if err != nil {
		log.Fatalf("error creating auth -- %s", err)
	}

	isAuthenticated := saAuth.IsAuthenticated()

	assert.Equal(t, true, isAuthenticated)
}

func TestChecksAuthNotDone(t *testing.T) {
	testAccountKeyPath := "../../testdata/test_account.json"
	accountKey, _ := NewServiceAccountKey(testAccountKeyPath)

	cmdLine := testutil.FakeCommandRunner("defaultOutput")
	cmdLine.FakeCommand(listAuthCmd, statusNotActiveJson)

	saAuth, err := newServiceAccountAuthWithRunner(cmdLine, accountKey)

	if err != nil {
		log.Fatalf("error creating auth -- %s", err)
	}

	isAuthenticated := saAuth.IsAuthenticated()

	assert.Equal(t, false, isAuthenticated)
}

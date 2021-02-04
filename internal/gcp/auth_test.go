package gcp

import (
	"fmt"
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

func TestAuthenticatesWhenValidSa(t *testing.T) {
	testAccountKeyPath := "../../testdata/test_account.json"
	accountKey, _ := NewServiceAccountKey(testAccountKeyPath)

	cmdLine := testutil.FakeCommandRunner("test-output")
	saAuth, err := newServiceAccountAuthWithRunner(cmdLine, accountKey)

	if err != nil {
		log.Fatalf("error creating auth -- %s", err)
	}

	saAuth.Authenticate()

	cmdLineRun := cmdLine.GetCmdLines()[0]
	expectedCmdLine := fmt.Sprintf("gcloud auth activate-service-account --key-file %s", testAccountKeyPath)
	assert.Equal(t, expectedCmdLine, cmdLineRun)
}

func TestAuthNotDone(t *testing.T) {
	testAccountKeyPath := "../../testdata/test_account.json"
	accountKey, _ := NewServiceAccountKey(testAccountKeyPath)
	cmdLine := testutil.FakeCommandRunner("defaultOutput")

	tests := map[string]struct {
		cmdLine string
		output  string
		auth    bool
	}{
		"auth done with active status":         {cmdLine: listAuthCmd, output: statusActiveJson, auth: true},
		"auth not done with non-active status": {cmdLine: listAuthCmd, output: statusNotActiveJson, auth: false},
		"auth not done with empty json":        {cmdLine: listAuthCmd, output: statusEmptyJson, auth: false},
		//"auth not done with empty string": {cmdLine: listAuthCmd, output: statusEmptyJson, auth: []bool{false}},
	}

	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {
			cmdLine.FakeCommand(tc.cmdLine, tc.output)
			saAuth, err := newServiceAccountAuthWithRunner(cmdLine, accountKey)

			if err != nil {
				t.Fatalf("error creating auth -- %s", err)
			}

			isAuthenticated := saAuth.IsAuthenticated()
			assert.Equal(t, tc.auth, isAuthenticated)
		})
	}
}

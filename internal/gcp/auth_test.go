package gcp

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/iac-io/myiac/testing"
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

	testAccountKeyJson = `
{
  "type": "service_account",
  "project_id": "testProject",
  "private_key_id": "45cc80629bd8fd2d018ff005f708bd3faa1bca45",
  "private_key": "-----BEGIN PRIVATE KEY-----\nIAABoGIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDIUXZKCIOzduR9\nMpXm8d2eEtD7dbMTaIWKj95MflJpgJzGuCWTg5g79U0gRARW+5h1bQFbIyIe4NYZ\nCNkItLLwYb34aAZc98yb7ulu0oZ41fZ0OOTOPPQZrX0CZaev2kIxYpzKNYqiT/b3\ngjmUTBhuY/sXJAJ/qzgRX0MoYXwDoZ9CeRF3NKwLdp4jNnXMRsGEXgVKEzm3hNmX\nP4lQKjj9H2+Gig93WCUZN6lbfdbp/13kMpzT0+cw6XTsdzLHYOotM3miadNCGIt0\n53oouHJeQ5PSuFzQszYbSe3tdm8BMYLdQiQog50quxmkU2VXxUJplLQhGaIXgQ3p\nVn/bG36dAgMBAAECggEAGbuDNUPuPSfU9rNAl+PyituSbncCa8gVvYS5MvzYM9bS\nbOGbbB1vuSYMBAzQvObBgTYhQjaba7mIrzsYfDqQMPpxV59vT9KCPXa9lF+laBDe\nQbRMSiUA22qSoDP0TE3+ik8HYp9poWuhx046fM8gpU+hIeodiw5w26Rv4VhSgLmw\nOF6yUx2po/cQdz0/+qScu82D3MnEtfdm8GWBSu5NU+lT4wvrCa8RxGoZ8cCsGDJn\n3qTn6f9qzSMgZcd4Sb2LGw5/UbpL43uUmWCRVMR/yCBFdYFG48rPzD8Nb8DZCjLD\nXOM3Crx20+vSI2/Itb34/Fju5v11whFdhhwBpRajgQKBgQD1WLQdklQs9FtboxfR\nv6930zgjfv1LBx9uR7qZy8f7F1fLBRW1BeVyNiGjRs3ONj0g059KhP6F60pM6iN5\nq28Zf0MpPzSpZ47J/emVkUIBGoy71tDWirAUpriUbTljzGAzRXSBySSr6U70cg+l\n7QE1DK+9vlW2BxgptWtX65JRgQKBgQDRBDgxvh82oYbJIbPyVL9Dak3jABsveXPX\n2/Y3TjVdz+0/pSb2KTmb0LNaYnvXv6+3xLS5KM4xy5pZyGNV8PXAByF6/Ya+1pvY\nnz1dLFjQzbp1rkPJWiTdsEtUM1mYivKTZC7KfX14bItlSBmkmJvBFs5BSpo9b5dz\njGIPWvPDHQKBgGn6gAsKC0xD3TavM3nJ+CylU2mZ0CXZlM0ZNNR8Pw0KH0U2FBNW\n0a7NDSivS/UYXr1QTE1vN1Z3tWeV9+71i48S9trZT5Ehh39fK8gMr9s0Mbht6VXT\nII47Gh4bNCAUxzU+ej4Zubp8lDtpDbNZthzJNxyaHAH9/IT/tbeLrW+BAoGASZSB\nr8ktNc8xItcRgOqiljnzB0l/SHwp8sCFcby/frH25CPgjmG+3QJgUR5AWJgrZLcD\no/cgd1kkkhzAE34LFTmtaJ2ddMsZ++067fTxozf5PvpE9LoeJkisjAyzqsanVIm9\nCx2YMO+NNu9lz5LFqfi8TTHVEHGbUFsIHj23eGUCgYBPIKpfee5ywRLpc5Rg2V6V\nY0wB4dmO1hYnY18XUm2xaWWEOfhiK9SDObUhmNWSRJ6dZCPR2shcYGa5Ew+LxGuS\nQV8eMQYw8PcwxFzNsU4Q4s2IucWGH5D0G7wFQys+ZqB1tYdvvMbQTOQ9X4R/Qc54\nWbyS1uABSbQxTCFNj0Fabc==\n-----END PRIVATE KEY-----\n",
  "client_email": "testAccount@testProject.iam.gserviceaccount.com",
  "client_id": "1000",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/testAccount%40testProject.iam.gserviceaccount.com"
}
`
	testAccountKeyPath = "/tmp/test_account.json"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func setup() {
	WriteSaKeyToString(testAccountKeyPath)
}

func cleanup() {
	_ = os.Remove(testAccountKeyPath)
}

func TestAuthenticatesWhenValidSa(t *testing.T) {
	accountKey, _ := NewServiceAccountKey("/tmp/test_account.json")

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

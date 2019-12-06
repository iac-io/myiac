package test
import (
	"testing"
	"github.com/dfernandezm/myiac/app/deploy"
)

const EXISTING_RELEASES_OUTPUT = `
{
	"Next": "",
	"Releases": [{
		"Name": "esteemed-peacock",
		"Revision": 2,
		"Updated": "Mon Dec  2 18:26:30 2019",
		"Status": "DEPLOYED",
		"Chart": "moneycolfrontend-1.0.0",
		"AppVersion": "0.1.0",
		"Namespace": "default"
	}, {
		"Name": "opining-frog",
		"Revision": 36,
		"Updated": "Fri Dec  6 13:41:17 2019",
		"Status": "DEPLOYED",
		"Chart": "traefik-1.78.4",
		"AppVersion": "1.7.14",
		"Namespace": "default"
	}, {
		"Name": "ponderous-lion",
		"Revision": 3,
		"Updated": "Mon Dec  2 18:26:30 2019",
		"Status": "DEPLOYED",
		"Chart": "moneycolserver-1.0.0",
		"AppVersion": "1.0.0",
		"Namespace": "default"
	}, {
		"Name": "solitary-ragdoll",
		"Revision": 2,
		"Updated": "Thu Dec  5 12:48:25 2019",
		"Status": "DEPLOYED",
		"Chart": "elasticsearch-1.0.0",
		"AppVersion": "6.5.0",
		"Namespace": "default"
	}]
}
`
//https://stackoverflow.com/questions/19167970/mock-functions-in-go
//TODO: this executes the real command, it should be mocked
//https://quii.gitbook.io/learn-go-with-tests/
// To run: go test github.com/dfernandezm/myiac/test
func TestReleaseDeployed(t *testing.T) {
    deployed := deploy.ReleaseDeployedForApp("traefik") != ""
    if !deployed {
       t.Errorf("The release is deployed was incorrect, got: %v, want: %v.", deployed, true)
    }
}
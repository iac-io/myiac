package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const CloudRunBasicYaml = `apiVersion: v1
kind: CloudRunDeployment
metadata:
 name: test-deployment
 namespace: default
 region: europe-west1
 annotations:
  custom.domain/provider: "cloudflare"
  custom.domain/dns: "test-deployment.moneycol.net"
spec:
 type: FullyManaged
 container:
  concurrency: 80
  cpu: 2
  port: 8080
  env:
   - name: aName
   - value: aValue
`

func TestMarshallsWithoutError(t *testing.T) {

	cloudRunDeploy, err := NewCloudRunDeployment(CloudRunBasicYaml)

	if err == nil {
		assert.Fail(t, "error marshalling yaml %s", CloudRunBasicYaml)
		panic(err)
	}

	if cloudRunDeploy == nil {
		assert.Fail(t, "cloudRun value is empty for yaml %s", CloudRunBasicYaml)
		panic("")
	}

	assert.Equal(t, "v1", cloudRunDeploy.cloudRunDeployment.APIVersion)
	assert.Equal(t, "CloudRunDeployment", cloudRunDeploy.cloudRunDeployment.Kind)
	assert.Equal(t, "test-deployment", cloudRunDeploy.cloudRunDeployment.Metadata.Name)
	assert.Equal(t, "default", cloudRunDeploy.cloudRunDeployment.Metadata.Namespace)
	assert.Equal(t, "europe-west1", cloudRunDeploy.cloudRunDeployment.Metadata.Region)
}

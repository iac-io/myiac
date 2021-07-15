package deploy

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// https://stackoverflow.com/questions/28682439/go-parse-yaml-file
type CloudRunDeployment struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Namespace   string `yaml:"namespace"`
		Region      string `yaml:"region"`
		Annotations struct {
			CustomDomainProvider string `yaml:"custom.domain/provider"`
			Domain               string `yaml:"custom.domain/value"`
			DnsProviderSecret    string `yaml:"custom.domain/secret"`
		} `yaml:"annotations"`
	} `yaml:"metadata"`
	Spec struct {
		Type           string `yaml:"type"`
		RequestTimeout int    `yaml:"requestTimeout"`
		Concurrency    string `yaml:"concurrency"`
		Container      struct {
			Cpu         int    `yaml:"cpu"`
			MemoryLimit string `yaml:"memoryLimit"`
			Port        int    `yaml:"port"`
			Command     int    `yaml:"command.omitempty"`
			Args        string `yaml:"args,omitempty"`
			Env         []struct {
				Name  string `yaml:"name"`
				Value string `yaml:"value"`
			} `yaml:"env,omitempty"`
		} `yaml:"container"`
	} `yaml:"spec"`
}

type cloudRunDeploy struct {
	cloudRunDeployment CloudRunDeployment
}

// NewCloudRunDeployment create a new CloudRunDeployment entity from its yaml representation
// returns a CloudRunDeploy struct or an error if it could not be marshalled correctly
func NewCloudRunDeployment(deploymentYaml string) (*cloudRunDeploy, error) {
	var cloudRunDeployment CloudRunDeployment
	cloudRunYamlBytes := []byte(deploymentYaml)
	err := yaml.Unmarshal(cloudRunYamlBytes, &cloudRunDeployment)

	if err != nil {
		return nil, errors.Errorf("could not marshal yaml %s", deploymentYaml)
	}

	return &cloudRunDeploy{cloudRunDeployment: cloudRunDeployment}, nil
}

// Input validation: https://link.medium.com/6fLKhIxQ5gb
func (crd cloudRunDeploy) Deploy() error {
	return nil
}

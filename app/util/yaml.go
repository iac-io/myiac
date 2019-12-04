package util

import (
	"fmt"
	"log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ExternalIpList struct {
	ExternalIPs []string
}

// Writes the array given to a yaml file at 'yamlFilePath'
func WriteAsYaml(yamlFilePath string, values []string) string {
	externalIpListS := ExternalIpList{ExternalIPs: values}
	ipsYaml, err := yaml.Marshal(externalIpListS)
	if err != nil {
		log.Fatalf("error marhsalling yaml: %v", err)
	}
	writeFile(yamlFilePath, string(ipsYaml))
	return yamlFilePath
}

func writeFile(filePath string, contents string) {
    data := []byte(contents)
    err := ioutil.WriteFile(filePath, data, 0755)
	check(err)
	fmt.Printf("Wrote file %s with content:\n %s\n",filePath, contents)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

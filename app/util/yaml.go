package util

import (
	"fmt"
	"log"
	"gopkg.in/yaml.v2"
)

// type ExternalIpList struct {
// 	IpList []string
// }
type ExternalIPs struct {
	
		IpList []string
	
}

type ExternalIpList struct {
	ExternalIPs []string
}

func BuildYamlList(values []string) {
	yamlList := ""
	for _, value := range values {
		yamlList += "- " + value + "\n"
	}
	//fmt.Print("list: %s\n", yamlList)
	externalIpList := make([]string,2)
	err := yaml.Unmarshal([]byte(yamlList), &externalIpList)
	if err != nil {
			log.Fatalf("error: %v", err)
	}

	//fmt.Printf("--- t:\n%v\n\n", externalIpList)
	//m := make(map[interface{}]interface{})
	//externalIpListS := ExternalIpList{ExternalIPs: ExternalIPs { IpList: values}}
	externalIpListS := ExternalIpList{ExternalIPs: values}
	d, err := yaml.Marshal(externalIpListS)
        if err != nil {
                log.Fatalf("error: %v", err)
        }
        fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
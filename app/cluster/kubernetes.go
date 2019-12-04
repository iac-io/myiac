package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

func GetInternalIpsForNodes() []string {
	json := executeGetIpsCmd()
	ips := getAllIps(json, true)
	fmt.Printf("Internal IPs for nodes in cluster are: %v\n", ips)
	return ips
}

func GetAllPublicIps() []string {
	json := executeGetIpsCmd()
	ips := getAllIps(json, false)
	fmt.Printf("Public IPs for nodes in cluster are: %v\n", ips)
	return ips
}

func executeGetIpsCmd() map[string]interface{} {
	argsArray := []string{"get", "nodes", "-o", "json"}
	cmd := commandline.New("kubectl", argsArray)
	cmd.SupressOutput = true
	cmdResult := cmd.Run()
	cmdOutput := cmdResult.Output()
	json := util.Parse(cmdOutput)
	return json
}

func getAllIps(json map[string]interface{}, internal bool) []string {
	indexOfAddress := 1
	if internal {
		indexOfAddress = 0
	} 
	allNodes := util.GetJsonArray(json, "items")
	var ips []string
	for _, node := range allNodes {
		status := util.GetJsonObject(node, "status")
		addresses := util.GetJsonArray(status, "addresses")
		ip := util.GetStringValue(addresses[indexOfAddress], "address") 
		ips = append(ips, ip)
	}
	return ips
}


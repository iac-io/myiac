package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

func GetInternalIpsForNodes() []string {
	argsArray := []string{"get", "nodes", "-o", "json"}
	cmd := commandline.New("kubectl", argsArray)
	cmd.SupressOutput = true
	cmdResult := cmd.Run()
	cmdOutput := cmdResult.Output()
	json := util.Parse(cmdOutput)
	ips := getAllInternalIps(json)
	fmt.Printf("IPs for nodes are: %v", ips)
	return ips
}

func getAllInternalIps(json map[string]interface{}) []string {
	allNodes := util.GetJsonArray(json, "items")
	var ips []string
	for _, node := range allNodes {
		status := util.GetJsonObject(node, "status")
		addresses := util.GetJsonArray(status, "addresses")
		ip := util.GetStringValue(addresses[0], "address")
		ips = append(ips, ip)
	}
	return ips
}
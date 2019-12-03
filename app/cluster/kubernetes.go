package cluster

import (
	//"fmt"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

func GetInternalIpsForNodes() {
	//getNodesIps := "get nodes -o json"
	argsArray := []string{"get", "nodes", "-o", "json"}
	cmd := commandline.New("kubectl", argsArray)
	cmdResult := cmd.Run()
	cmdOutput := cmdResult.Output()
	//fmt.Printf("Nodes Json is: %s\n", cmdOutput)
	util.JsonAsMap(cmdOutput, "")
}
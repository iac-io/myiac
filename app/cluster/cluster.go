package cluster

import (
	"fmt"
	"strings"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

func InstallHelm() {
	fmt.Println("Installing Helm in newly created cluster")
	
	//TODO: this should be a configurable path
	// the directory where the binary is being executed from
	// The script itself should be inlined, read an executed instead of 
	// bundling a file with the binary
	currentDir := util.CurrentExecutableDir() 
	helperScriptsLocation := currentDir + "/helperScripts"
	fmt.Printf("Helper scripts path is %s", helperScriptsLocation)

	action := "./install-helm.sh"
	cmd := commandline.NewWithWorkingDir(action, []string{}, helperScriptsLocation)
	cmd.Run()
}

//unused: delete
func labelNodes(nodeType string) {
	//slice vs array: https://blog.golang.org/go-slices-usage-and-internals
	var nodeNames []string
	var label string
	nodeNamesEs := []string{"gke-moneycol-main-elasticsearch-pool-b8711571-k359"}
	nodeNamesApps := []string{"gke-moneycol-main-main-pool-ac0c4442-57ff",
		"gke-moneycol-main-main-pool-ac0c4442-pq57",
		"gke-moneycol-main-main-pool-ac0c4442-q1t7"}

	if nodeType == "elasticsearch" {
		nodeNames = nodeNamesEs
		label = "type=elasticsearch"
	} else if nodeType == "applications" {
		nodeNames = nodeNamesApps
		label = "type=applications"
	}

	labelCmdTpl := "label nodes %s %s --overwrite\n"

	//note: range (like everything in go) copies by value the slice
	for _, nodeName := range nodeNames {
		argsStr := fmt.Sprintf(labelCmdTpl, nodeName, label)
		fmt.Printf("Labelling args: %s", argsStr)
		argsArray := strings.Fields(argsStr)
		c := commandline.New("kubectl", argsArray)
		c.Run()
	}
}
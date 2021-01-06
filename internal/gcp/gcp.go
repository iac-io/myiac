package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/iac-io/myiac/internal/commandline"
	"github.com/iac-io/myiac/internal/util"
	"log"
	"strconv"
	"strings"
)

type nodePoolList []map[string]interface{}

// Deprecated
// see provider.go / setup_environment.go
func SetupEnvironment(projectId string) {
	keyLocation := util.GetHomeDir() + fmt.Sprintf("/%s_account.json", projectId)
	baseArgs := "auth activate-service-account --key-file %s"
	var argsArray []string = util.StringTemplateToArgsArray(baseArgs, keyLocation)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
}

func ConfigureDocker() {
	action := "auth configure-docker"
	cmdTpl := "%s"
	argsArray := util.StringTemplateToArgsArray(cmdTpl, action)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
}

// Deprecated
func SetupKubernetes(project string, zone string, environment string) {
	action := "container clusters get-credentials"
	// gcloud container clusters get-credentials [cluster-name]
	cmdTpl := "%s %s --zone %s --project %s"
	clusterName := fmt.Sprintf("%s-%s", project, environment)

	argsArray := util.StringTemplateToArgsArray(cmdTpl, action, clusterName, zone, project)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
	fmt.Println("Kubernetes setup completed")
}

// https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
func parseNodePoolList(jsonString string) nodePoolList {
	var listNodePools nodePoolList
	
	// If there is no releases, a single space is returned
	if jsonString == "" || len(strings.TrimSpace(jsonString)) == 0 {
		// empty releases list
		fmt.Printf("Empty list of releases found")
	} else {
		jsonData := []byte(jsonString)
		err := json.Unmarshal(jsonData, &listNodePools)
		if err != nil {
			log.Fatalf("Error parsing json to struct %v", err)
		}
	}
	return listNodePools
}

// ListClusterNodePools list the node pools in the GKE cluster
func ListClusterNodePools(project string, zone string, environment string) nodePoolList {
	cmdTpl := "container node-pools list --cluster %s --zone %s --project %s --format json"
	clusterName := fmt.Sprintf("%s-%s", project, environment)
	argsArray := util.StringTemplateToArgsArray(cmdTpl, clusterName, zone, project)
	cmd := commandline.New("gcloud", argsArray)
	cmd.Run()
	nodePools := parseNodePoolList(cmd.Output())
	fmt.Printf("Node Pools in cluster %s: %v\n", clusterName, nodePools)
	return nodePools	
}

// ResizeCluster change the size of GKE cluster node pools to the provided values
//
// The resize command needs to be executed once per node pool:
// 
// gcloud container clusters resize NAME (--num-nodes=NUM_NODES | --size=NUM_NODES) [--async]
// [--node-pool=NODE_POOL] [--region=REGION | --zone=ZONE, -z ZONE]
func ResizeCluster(project string, zone string, environment string, targetSize int) {
	clusterName := fmt.Sprintf("%s-%s", project, environment)
	action := "container clusters resize"

	nodePools := ListClusterNodePools(project, zone, environment)
	nodePoolResizeTpl := "--node-pool %s --num-nodes %s"
	
	for _,np := range nodePools {
		nodePoolName := np["name"]
		nodePoolNameStr, _ := nodePoolName.(string)
		fmt.Printf("Resizing node pool %s to %d\n", nodePoolName, targetSize)
		targetSizeStr := strconv.Itoa(targetSize)
		nodePoolResizeArray := util.StringTemplateToArgsArray(nodePoolResizeTpl, nodePoolNameStr, targetSizeStr)
		nodePoolResizePart := strings.Join(nodePoolResizeArray, " ")
		
		cmdTpl := "%s %s %s --zone %s --project %s -q"
		argsArray := util.StringTemplateToArgsArray(cmdTpl, action, clusterName, nodePoolResizePart, zone, project)
		cmd := commandline.New("gcloud", argsArray)
		cmd.Run()
		fmt.Printf("Node Pool %s resized", nodePoolName)	
	}
	fmt.Println("Cluster resized")
}

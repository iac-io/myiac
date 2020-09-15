package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dfernandezm/myiac/internal/commandline"
	"github.com/dfernandezm/myiac/internal/util"
	iam "google.golang.org/api/iam/v1"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type nodePoolList []map[string]interface{}


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
//	gcloud container clusters resize NAME (--num-nodes=NUM_NODES | --size=NUM_NODES) [--async]
// 	[--node-pool=NODE_POOL] [--region=REGION | --zone=ZONE, -z ZONE]
func ResizeCluster(project string, zone string, environment string, targetSize int) {
	clusterName := fmt.Sprintf("%s-%s", project, environment)
	action := "container clusters resize"

	nodePools := ListClusterNodePools(project, zone, environment)
	nodePoolResizeTpl := "--node-pool %s --num-nodes %s"
	
	for _,np := range nodePools {
		nodePoolName := np["name"]
		nodePoolNameStr, _ := nodePoolName.(string)
		fmt.Printf("Resizing node pool %s to %d", nodePoolName, targetSize)
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


// CreateKey creates a service account key for the given service account email
func CreateKey(w io.Writer, serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := &iam.CreateServiceAccountKeyRequest{}

	fmt.Printf("Creating the key\n")

	key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}
	_, _ = fmt.Fprintf(w, "Created key: %v", key.Name)
	fmt.Printf("Created the key %v\n", key)
	return key, nil
}

func ListKeys(serviceAccountEmail string) ([]*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	response, err := service.Projects.ServiceAccounts.Keys.List(resource).Do()

	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.List: %v", err)
	}

	return response.Keys, nil
}

func KeyForServiceAccount(saEmail string, recreateKey bool) (string, error) {
	keys, _ := ListKeys(saEmail)
	var key *iam.ServiceAccountKey = nil
	var jsonKey string

	if len(keys) > 0 && !recreateKey {
		fmt.Printf("Using existing key for SA: %s\n", saEmail)
		key = keys[0]
		keyString, err := findKeyInObjectStorage(saEmail)
		if err != nil {
			return "", fmt.Errorf("err finding key %v", err)
		}
		jsonKey = keyString
	} else {
		newKey, err := CreateKey(os.Stdout, saEmail)
		if err != nil {
			fmt.Printf("Error creating key %v", err)
			return "", err
		}
		key = newKey
		privateKeyData := key.PrivateKeyData
		jsonKeyString := util.Base64Decode(privateKeyData)

		fmt.Println("----- BEGIN Service account JSON key -------")
		fmt.Println(jsonKeyString)
		fmt.Println("----- END Service account JSON key ---------")

		writeErr := writeKeyToObjectStorage(saEmail, jsonKeyString)
		if writeErr != nil {
			return "", fmt.Errorf("error writing key to storage %v", writeErr)
		}

		jsonKey = jsonKeyString
	}

	return jsonKey, nil
}

func writeKeyToObjectStorage(saEmail string, jsonKey string) error {
	ctx := context.Background()
	bucketName := "moneycol-keys"
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(saEmail + ".json")
	// Write something to obj.
	// w implements io.Writer.
	w := obj.NewWriter(ctx)

	// Write some text to obj. This will either create the object or overwrite whatever is there already.
	if _, err := fmt.Fprintf(w, "%s", jsonKey); err != nil {
		return err
	}

	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func findKeyInObjectStorage(saEmail string) (string, error) {
	ctx := context.Background()

	//TODO: global parameter
	bucketName := "moneycol-keys"
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(saEmail + ".json")

	//TODO: should list the bucket first, see:
	// https://godoc.org/cloud.google.com/go/storage
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("getting object error %v", err)
	}

	defer r.Close()

	buf := new(strings.Builder)

	if _, copyErr := io.Copy(buf, r); copyErr != nil {
		return "", fmt.Errorf("error copying to stdout %v", copyErr)
	}

	keyString := buf.String()
	fmt.Printf("Found key for account %s in Storage\n", saEmail)
	fmt.Printf("Found key %s\n", keyString)
	return keyString, nil
}

package util

import (
	"encoding/json"
	"fmt"
)

type JSON map[string]interface{}

//https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
func JsonAsMap(jsonString string, field string) {
	var objmap map[string]*json.RawMessage
	data := []byte(jsonString)

	err := json.Unmarshal(data, &objmap)

	if err == nil {
		//TODO: throw error
	}
	
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonMap)

	allNodes := getJsonArray(jsonMap, "items")
	//fmt.Printf("Nodes: %v",allNodes)
	var ips []string
	for _, node := range allNodes {
		status := getJsonObject(node, "status")
		addresses := getJsonArray(status, "addresses")
		ip := getStringValue(addresses[0], "address")
		ips = append(ips, ip)
	}

	fmt.Printf("internal IPs %v\n", ips)
}

func getJsonObject(jsonMap interface{}, key string) interface{} {
	return jsonMap.(map[string]interface{})[key].(interface{})
}

func getJsonArray(jsonMap interface{}, arrayKey string) []interface{} {
	array := jsonMap.(map[string]interface{})[arrayKey]
	return array.([]interface{})
}

func getStringValue(jsonMap interface{}, key string) string {
	return jsonMap.(map[string]interface{})[key].(string)
}

func getString(value interface{}) string {
	return value.(string)
}

func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

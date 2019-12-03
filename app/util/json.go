package util

import (
	"encoding/json"
	"fmt"
	"log"
)

type JSON map[string]interface{}

//https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
func Parse(jsonString string) map[string]interface{} {
	var objmap map[string]*json.RawMessage
	data := []byte(jsonString)

	err := json.Unmarshal(data, &objmap)

	if err != nil {
		log.Fatalf("Error unmarshalling json: %v", err)
	}

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonMap)

	return jsonMap
}

func GetJsonObject(jsonMap interface{}, key string) interface{} {
	return jsonMap.(map[string]interface{})[key].(interface{})
}

func GetJsonArray(jsonMap interface{}, arrayKey string) []interface{} {
	array := jsonMap.(map[string]interface{})[arrayKey]
	return array.([]interface{})
}

func GetStringValue(jsonMap interface{}, key string) string {
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

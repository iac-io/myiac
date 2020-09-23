package util

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// StringTemplateToArgsArray spreads an array of values on top of a templated string (%s,%s,...).
// It receives a templated string (%s, %d, ...) and a spread of values ['a','b','c'] to generate a string
func StringTemplateToArgsArray(templatedArgs string, values ...string) []string {
	baseArgsTmpl := fmt.Sprintf(templatedArgs, stringToInterfaceArrays(values)...)
	return strings.Fields(baseArgsTmpl)
}

func LookupPathFor(executable string) (string, error) {
	path, err := exec.LookPath(".")
	return path, err
}

func CurrentExecutableDir() string {
	
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting the directory of this executable %s", err)
	}

	location, err := filepath.Abs(filepath.Dir(executable))

	if (err != nil) {
		log.Fatalf("Error getting the directory of this executable %s", err)
	}

	fmt.Printf("This executable full path directory is: %s\n", location)
	return location
}

func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// stringToInterfaceArrays converts a spread of array (arguments) of type string into an interface{} array.
// This way it can be spread into Sprinf family functions (they receive ...interface{})
// See: https://golang.org/doc/faq#convert_slice_of_interface and https://github.com/golang/go/issues/15037
func stringToInterfaceArrays(arr []string) []interface{} {
	s := make([]interface{}, len(arr))
	for i, v := range arr {
		s[i] = v
	}
	return s
}

func Base64Decode(toDecode string) string {
	decodedBytes, _ := b64.StdEncoding.DecodeString(toDecode)
	decodedString := string(decodedBytes)
	return decodedString
}

// Writes content to a filePath, overriding the current contents
// Returns: error when cannot create the file or write content
func WriteStringToFile(content string, filePath string) error {
	fmt.Printf("Creating file\n")

	fo, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s %v\n", filePath, err)
		return err
	}

	bytes, err := fmt.Fprintf(fo, "%s", content)
	if err != nil {
		fmt.Printf("Error creating file %s %v\n", filePath, err)
		return err
	}

	fmt.Printf("Written %d bytes\n", bytes)
	return nil
}

func ArrayContains(array []string, value string) bool {
	for _, val:= range array {
		if val == value {
			return true
		}
	}
	return false
}

func ReadFileToString(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file %v", err)
	}
	return string(bytes), nil
}

func ReadFileToBytes(filename string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %v", err)
	}
	return bytes, nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}



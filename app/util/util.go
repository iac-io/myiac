package util

import (
	"path/filepath"
	"fmt"
	"log"
	"os"
	"os/user"
	"os/exec"
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

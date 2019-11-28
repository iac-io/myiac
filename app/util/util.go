package util

import (
	"fmt"
	"log"
	"os/user"
	"strings"
)

// StringTemplateToArgsArray spreads an array of values on top of a templated string (%s,%s,...).
// It receives a templated string (%s, %d, ...) and a spread of values ['a','b','c'] to generate a string
func StringTemplateToArgsArray(templatedArgs string, values ...string) []string {
	baseArgsTmpl := fmt.Sprintf(templatedArgs, stringToInterfaceArrays(values)...)
	return strings.Fields(baseArgsTmpl)
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

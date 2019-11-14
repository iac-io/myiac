package util

import (
	"fmt"
	"log"
	"os/user"
	"strings"
)

func StringTemplateToArgsArray(templatedArgs string, values ...string) []string {
	baseArgsTmpl := fmt.Sprintf(templatedArgs, values[0]) //TODO: spread operator how to use
	return strings.Fields(baseArgsTmpl)
}

func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

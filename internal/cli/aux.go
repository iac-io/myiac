package cli

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"strings"
)

func validateBaseFlags(ctx *cli.Context) error {
	validateStringFlagPresence("project", ctx)
	return nil
}

func validateNodePoolsSize(ctx *cli.Context) error {
	fmt.Printf("Validating flag nodePoolsSize\n")
	nodePoolsSizeValue := ctx.String("nodePoolsSize")
	fmt.Printf("Flag nodePoolsSize read as %s\n", nodePoolsSizeValue)

	val, err := strconv.Atoi(nodePoolsSizeValue)

	if err != nil {
		log.Printf("Error converting %v\n", err)
		logErrorAndExit("Invalid nodePoolsSize: " + nodePoolsSizeValue)
	}

	if val >= 0 {
		// Valid, it's greater than 0
		fmt.Printf("Valid nodePoolSize %s", nodePoolsSizeValue)
		return nil
	}

	if val < 0 {
		logErrorAndExit(fmt.Sprintf("Invalid nodePoolsSize: %d", val))
	}

	if len(nodePoolsSizeValue) == 0 {
		logErrorAndExit("Invalid nodePoolsSize: " + nodePoolsSizeValue)
	}

	return nil
}

func logErrorAndExit(errorMsg string) {
	err := cli.NewExitError(errorMsg, -1)
	if err != nil {
		log.Fatalf(errorMsg, err)
	}
}

func validateStringFlagPresence(flagName string, ctx *cli.Context) string {
	log.Printf("Validating flag %s\n", flagName)
	flag := ctx.String(flagName)
	if flag == "" {
		log.Fatalf("%s parameter not provided\n", flagName)
	}
	log.Printf("Read flag %s as %s\n", flagName, flag)
	return flag
}

func readPropertiesToMap(properties string) map[string]string {
	propertiesMap := make(map[string]string)
	if len(properties) > 0 {
		fmt.Printf("Properties line is %s\n", properties)
		keyValues := strings.Split(properties, ",")
		for _, s := range keyValues {
			keyValue := strings.Split(s, "=")
			key := keyValue[0]
			value := keyValue[1]
			fmt.Printf("Read property: %s -> %s\n", key, value)
			propertiesMap[key] = value
		}
	}
	return propertiesMap
}

func setHelmChartsPathVar(k string) {
	log.Printf("Setting Helm charts path to: %s", k)
	err := os.Setenv("CHARTS_PATH", k)
	if err != nil {
		log.Fatalf("Could not setup Helm charts path: %s", k)
	}
}

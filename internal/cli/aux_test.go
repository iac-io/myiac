package cli

import (
	"log"
	"os"
	"testing"
)

func TestSetHelmChartsPathVar(t *testing.T) {
	setHelmChartsPathVar("/test")
	if os.Getenv("CHARTS_PATH") != "/test" {
		log.Panic("HELM charts path failed.")
	}
}

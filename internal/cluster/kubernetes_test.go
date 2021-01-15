package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIpsAsHelmParams(t *testing.T) {
	ipArray := []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}

	helmParams := getNodesInternalIpsAsHelmParams(ipArray)

	expectedOutputIps := "{1.1.1.1\\,2.2.2.2\\,3.3.3.3}"

	assert.True(t, mapHasKey(helmParams, externalIpsKeyName))
	assert.Equal(t, helmParams[externalIpsKeyName], expectedOutputIps)
}

func mapHasKey(aMap map[string]string, key string) bool {
	if _, ok := aMap[key]; ok {
		return true
	}
	return false
}

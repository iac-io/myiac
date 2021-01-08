package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGCPDNSService(t *testing.T) {
	gcpClient := NewGoogleCloudDNSService("moneycol", "test-zone")
	assert.Equal(t, gcpClient.project, "moneycol")
	assert.Equal(t, gcpClient.zone, "test-zone")
	assert.NotNil(t, gcpClient.service)
}

func TestGetDnsRecord(t *testing.T) {
	gcpClient := NewGoogleCloudDNSService("moneycol", "test-zone")
	gcpClient.UpsertDNSRecord("A", "dev-test.moneycol-test.ml", "34.77.93.11")
	result := gcpClient.GetDNSRecordByName("A", "dev-test.moneycol-test.ml")

	assert.Len(t, result, 1)
	assert.Equal(t, result[0].Name, "dev-test.moneycol-test.ml.")
}

func TestChangeDnsRecord(t *testing.T) {
	// money-zone-free is pre-created in GCP
	gcpClient := NewGoogleCloudDNSService("moneycol", "test-zone")

	gcpClient.UpsertDNSRecord("A", "devtest.moneycol-test.ml", "34.77.93.11")
	result := gcpClient.GetDNSRecordByName("A", "devtest.moneycol-test.ml")
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].Rrdatas[0], "34.77.93.11")

	gcpClient.UpsertDNSRecord("A", "devtest.moneycol-test.ml", "34.77.93.10")
	result = gcpClient.GetDNSRecordByName("A", "devtest.moneycol-test.ml")
	assert.Len(t, result, 1)
	assert.Equal(t, result[0].Rrdatas[0], "34.77.93.10")
}

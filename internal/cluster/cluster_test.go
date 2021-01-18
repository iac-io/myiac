package cluster

import (
	"testing"
)

func TestValidateTFVars(t *testing.T) {
	got, got1 := ValidateTFVars("/home/app/terraform")
	if got != "/home/app/terraform" {
		t.Errorf("ValidateTFVars() got = %v, want %v", got, "/home/app/terraform")
	}
	if got1 != "/home/app/terraform/cluster.tfvars" {
		t.Errorf("ValidateTFVars() got1 = %v, want %v", got1, "/home/app/terraform/cluster.tfvars")
	}
}

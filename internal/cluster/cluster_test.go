package cluster

import (
	"testing"

	"github.com/iac-io/myiac/internal/util"
)

func TestValidateTFVars(t *testing.T) {
	got, got1 := ValidateTFVars("/home/app/terraform")
	if got != util.CurrentExecutableDir()+"/internal/terraform/cluster" {
		t.Errorf("ValidateTFVars() got = %v, want %v", got, util.CurrentExecutableDir()+"/internal/terraform/cluster")
	}
	if got1 != util.CurrentExecutableDir()+"/internal/terraform/cluster/terraform.tfvars" {
		t.Errorf("ValidateTFVars() got1 = %v, want %v", got1, util.CurrentExecutableDir()+"/internal/terraform/cluster")
	}
}

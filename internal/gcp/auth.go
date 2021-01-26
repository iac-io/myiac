package gcp

import "github.com/iac-io/myiac/internal/commandline"

type ServiceAccountAuth struct {
	serviceAccountKey ServiceAccountKey
	commandRunner     commandline.CommandRunner
}

func NewServiceAccountAuth(keyLocation string) {

}

func newServiceAccountAuthWithRunner(commandRunner commandline.CommandRunner, keyLocation string) {

}

func Authenticate() {

}

func IsAuthenticated() {

}

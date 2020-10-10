package deploy

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/commandline"
	"github.com/dfernandezm/myiac/internal/util"
)

func ApplyDnsIpChange(tfFileLocation string, ip string) {
	// do no use single quotes here for the var (i.e. -var 'foo=bar') as it fails to execute
	inlinedVar := fmt.Sprintf("-var dev_ip=%s",ip) 

	argsArray := util.StringTemplateToArgsArray("%s %s", "plan", inlinedVar)
	cmd := commandline.NewWithWorkingDir("terraform", argsArray, tfFileLocation)
	cmd.Run()

	argsArray = util.StringTemplateToArgsArray("%s %s %s", "apply", inlinedVar, "-auto-approve")
	cmd = commandline.NewWithWorkingDir("terraform", argsArray, tfFileLocation)
	cmd.Run()

	fmt.Printf("Applied change of IP for DNS\n")
}

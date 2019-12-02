package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

func InstallHelm() {
	fmt.Println("Installing Helm in newly created cluster")
	
	currentDir := util.CurrentExecutableDir()
	helperScriptsLocation := currentDir + "/.helperScripts"
	fmt.Printf("Helper scripts path is %s", helperScriptsLocation)

	action := "./install-helm.sh"
	cmd := commandline.NewWithWorkingDir(action, []string{}, helperScriptsLocation)
	cmd.Run()
}

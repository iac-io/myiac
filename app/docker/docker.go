package docker

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/util"
	"github.com/dfernandezm/myiac/app/commandline"
)

type DockerImageDefinition struct {
	projectId string
	projectRepository string
	imageIdToTag string
	commit string
	containerName string
	version string
}

func TagDockerImage(dockerImageDef *DockerImageDefinition) {
	tag := fmt.Sprintf("%s-%s", dockerImageDef.version, dockerImageDef.commit)
	containerFullName := fmt.Sprintf("%s/%s/%s:%s", dockerImageDef.projectRepository, dockerImageDef.projectId, dockerImageDef.containerName, tag)
	fmt.Printf("The image tag to push is: %s\n", containerFullName)
	dockerTagCmdPart := "tag"
	dockerCmd := []string{dockerTagCmdPart, dockerImageDef.imageIdToTag, containerFullName}
	argsArray := util.StringTemplateToArgsArray("%s %s %s", dockerCmd)
	cmd := commandline.New("docker", argsArray)
	cmd.Run()
	
	//runtime.SetDockerImage(containerFullName)
	fmt.Printf("Docker image has been tagged with: %s\n", containerFullName)
}


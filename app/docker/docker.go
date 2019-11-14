package docker

import (
	"fmt"
	"github.com/dfernandezm/myiac/app/util"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/properties"
)

const TAG_CMD_PART = "tag"

type DockerImageDefinition struct {
	projectId string
	projectRepositoryUrl string
	imageIdToTag string
	commitHash string
	imageRepoName string
	version string
}

func TagDockerImage(dockerImg *DockerImageDefinition, runtimeProps *properties.RuntimeProperties) {
	tag := fmt.Sprintf("%s-%s", dockerImg.version, dockerImg.commitHash)
	fullImageDefinition := fmt.Sprintf("%s/%s/%s:%s", dockerImg.projectRepositoryUrl, dockerImg.projectId, dockerImg.imageRepoName, tag)
	fmt.Printf("The image tag to push is: %s\n", fullImageDefinition)

	runDockerTagCmd(TAG_CMD_PART, dockerImg.imageIdToTag, fullImageDefinition)

	runtimeProps.SetDockerImage(fullImageDefinition)
}

func runDockerTagCmd(tagCmdPart string, imageId string, fullImageDefinition string) {
	argsArray := util.StringTemplateToArgsArray("%s %s %s", tagCmdPart, imageId, fullImageDefinition)
	cmd := commandline.New("docker", argsArray)
	cmd.Run()
	fmt.Printf("Docker image has been tagged with: %s\n", fullImageDefinition)
}


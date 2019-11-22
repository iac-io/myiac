package docker

import (
	"fmt"

	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/properties"
	"github.com/dfernandezm/myiac/app/util"
)

const TAG_CMD_PART = "tag"
const BUILD_IMAGE_PART = "build"

type dockerImageDefinition struct {
	projectId            string
	projectRepositoryUrl string
	imageIdToTag         string
	commitHash           string
	imageRepoName        string
	version              string
}

func NewImageDefinition(projectRepoUrl string, projectId string, imageId string,
	commitHash string, imageRepoName string, version string) dockerImageDefinition {
	//TODO: pass in all params
	return dockerImageDefinition{}
}

func (di *dockerImageDefinition) TagImage(runtimeProps *properties.RuntimeProperties) {
	fullImageDefinition := generateTag(di.projectRepositoryUrl, di.projectId, di.imageRepoName, di.version, di.commitHash)
	runDockerTagCmd(TAG_CMD_PART, di.imageIdToTag, fullImageDefinition)
	runtimeProps.SetDockerImage(fullImageDefinition)
}

func generateTag(projectRepoUrl string, projectId string, imageRepoName string, version string, commitHash string) string {
	tag := fmt.Sprintf("%s-%s", version, commitHash)
	fullImageDefinition := fmt.Sprintf("%s/%s/%s:%s", projectRepoUrl, projectId, imageRepoName, tag)
	fmt.Printf("The image tag to push is: %s\n", fullImageDefinition)
	return fullImageDefinition
}

func runDockerTagCmd(tagCmdPart string, imageId string, fullImageDefinition string) {
	argsArray := util.StringTemplateToArgsArray("%s %s %s", tagCmdPart, imageId, fullImageDefinition)
	cmd := commandline.New("docker", argsArray)
	cmd.Run()
	fmt.Printf("Docker image has been tagged with: %s\n", fullImageDefinition)
}

type dockerBuildImageDefinition struct {
	buildPath string
	tag       string
}

func BuildImage(buildPath string, tag string) {

}

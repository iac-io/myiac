package docker

import (
	"fmt"

	"github.com/dfernandezm/myiac/app/commandline"
	p "github.com/dfernandezm/myiac/app/properties"
	"github.com/dfernandezm/myiac/app/util"
)

const TAG_CMD_PART = "tag"
const BUILD_IMAGE_PART = "build"
const PUSH_IMAGE_PART = "push"

func TagImage(runtimeProperties *p.RuntimeProperties, dockerProps *p.DockerProperties, imageId string, commitHash string, imageRepo string, appVersion string) {
	fullImageDefinition := generateTag(dockerProps.ProjectRepoUrl, dockerProps.ProjectId, imageRepo, appVersion, commitHash)
	fmt.Printf("Full image definition is %s\n", fullImageDefinition)
	fmt.Printf("Image id to tag is %s\n", imageId)
	runDockerTagCmd(TAG_CMD_PART, imageId, fullImageDefinition)
	runtimeProperties.SetDockerImage(fullImageDefinition)
}

func PushImage(runtime *p.RuntimeProperties) {
	dockerImage := runtime.GetDockerImage()
	fmt.Printf("Pushing previously built docker image: %s\n", dockerImage)
	runDockerPushCmd(dockerImage)
}

func BuildImage(buildPath string) {

}

func generateTag(projectRepoUrl string, projectId string, imageRepoName string, version string, commitHash string) string {
	tag := fmt.Sprintf("%s-%s", version, commitHash)
	fullImageDefinition := fmt.Sprintf("%s/%s/%s:%s", projectRepoUrl, projectId, imageRepoName, tag)
	fmt.Printf("The image to push is: %s\n", fullImageDefinition)
	return fullImageDefinition
}

func runDockerTagCmd(tagCmdPart string, imageId string, fullImageDefinition string) {
	argsArray := util.StringTemplateToArgsArray("%s %s %s", tagCmdPart, imageId, fullImageDefinition)
	cmd := commandline.New("docker", argsArray)
	cmd.Run()
	fmt.Printf("Docker image has been tagged with: %s\n", fullImageDefinition)
}

func runDockerPushCmd(imageToPush string) {
	argsArray := util.StringTemplateToArgsArray("%s %s", PUSH_IMAGE_PART, imageToPush)
	cmd := commandline.New("docker", argsArray)
	cmd.Run()
	fmt.Printf("Docker image %s has been pushed successfully\n", cmd)
}

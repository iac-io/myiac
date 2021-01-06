package docker

import (
	"fmt"

	"github.com/iac-io/myiac/internal/commandline"
	p "github.com/iac-io/myiac/internal/properties"
	"github.com/iac-io/myiac/internal/util"
)

const TAG_CMD_PART = "tag"
const BUILD_IMAGE_PART = "build"
const PUSH_IMAGE_PART = "push"

func TagImage(runtimeProperties *p.RuntimeProperties, dockerProps *p.DockerProperties, imageId string, commitHash string, imageRepo string, appVersion string) {
	fullImageDefinition := GenerateTag(dockerProps.ProjectRepoUrl, dockerProps.ProjectId, imageRepo, appVersion, commitHash)
	fmt.Printf("Full image definition is %s\n", fullImageDefinition)
	fmt.Printf("Image id to tag is %s\n", imageId)
	runDockerTagCmd(TAG_CMD_PART, imageId, fullImageDefinition)
	runtimeProperties.SetDockerImage(fullImageDefinition)
}

func PushImage(runtime *p.RuntimeProperties) {
	dockerImage := runtime.GetDockerImage()
	fmt.Printf("Pushing previously built docker image: %s\n", dockerImage)
	runDockerPushCmd(dockerImage)
	runtime.SetDockerImage("")
}

func BuildImage(runtime *p.RuntimeProperties, buildPath string, dockerProps *p.DockerProperties, commitHash string, imageRepo string, appVersion string) {
	fullImageDefinition := GenerateTag(dockerProps.ProjectRepoUrl, dockerProps.ProjectId, imageRepo, appVersion, commitHash)
	fmt.Printf("Full image definition to build is %s\n", fullImageDefinition)
	runDockerBuildCmd("build", buildPath, fullImageDefinition)
	runtime.SetDockerImage(fullImageDefinition)
}

func GenerateTag(projectRepoUrl string, projectId string, imageRepoName string, version string, commitHash string) string {
	tag := fmt.Sprintf("%s-%s", version, commitHash)
	fullImageDefinition := fmt.Sprintf("%s/%s/%s:%s", projectRepoUrl, projectId, imageRepoName, tag)
	fmt.Printf("The image to push is: %s\n", fullImageDefinition)
	return fullImageDefinition
}

func runDockerBuildCmd(buildCmdPart string, buildPath string, fullImageDefinition string) {
	argsArray := util.StringTemplateToArgsArray("%s %s -t %s", buildCmdPart, buildPath, fullImageDefinition)
	cmd := commandline.New("docker", argsArray)
	cmd.Run()
	fmt.Printf("Docker image has been built: %s\n", fullImageDefinition)
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
	fmt.Printf("Docker image %s has been pushed successfully\n", imageToPush)
}

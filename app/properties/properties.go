package properties

type RuntimeProperties struct {
	dockerImage string //TODO: should be a map (appname, image)
}

func NewRuntime() RuntimeProperties {
	return RuntimeProperties{dockerImage: ""}
}

func (r *RuntimeProperties) SetDockerImage(dockerImage string) {
	r.dockerImage = dockerImage
}

func (r *RuntimeProperties) GetDockerImage() string {
	return r.dockerImage
}

type DockerProperties struct {
	ProjectId string
	ProjectRepoUrl string
}


package properties

type RuntimeProperties struct {
	dockerImage string //TODO: should be a map (appname, image)
}

func (r *RuntimeProperties) SetDockerImage(dockerImage string) {
	r.dockerImage = dockerImage
}

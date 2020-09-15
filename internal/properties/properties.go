package properties

type propertyStore struct {
	runtime RuntimeProperties
	helmProperties HelmProperties
}

var PropertyStore = propertyStore{}

func (ps *propertyStore) init() *propertyStore {
	//TODO: loadAllProperties
	return ps
}

func (ps *propertyStore) Get() *propertyStore {
	// return the created 
	return ps
}

type HelmProperties struct {
	appName string
}

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




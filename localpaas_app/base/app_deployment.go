package base

type BuildTool string

const (
	BuildToolDocker   BuildTool = "docker"
	BuildToolNixpacks BuildTool = "nixpacks"
)

var (
	AllBuildTools = []BuildTool{BuildToolDocker, BuildToolNixpacks}
)

type DeploymentSource string

const (
	DeploymentSourceImage   DeploymentSource = `image`
	DeploymentSourceRepo    DeploymentSource = "repo"
	DeploymentSourceTarball DeploymentSource = "tarball"
)

var (
	AllDeploymentSources = []DeploymentSource{DeploymentSourceImage, DeploymentSourceRepo, DeploymentSourceTarball}
)

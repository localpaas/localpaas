package base

type BuildTool string

const (
	BuildToolDocker   BuildTool = "docker"
	BuildToolNixpacks BuildTool = "nixpacks"
)

var (
	AllBuildTools = []BuildTool{BuildToolDocker, BuildToolNixpacks}
)

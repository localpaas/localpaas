package base

type BuildTool string

const (
	BuildToolDockerfile BuildTool = "dockerfile"
	BuildToolNixpacks   BuildTool = "nixpacks"
)

var (
	AllBuildTools = []BuildTool{BuildToolDockerfile, BuildToolNixpacks}
)

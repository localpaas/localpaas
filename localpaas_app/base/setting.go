package base

type SettingTargetType string

const (
	SettingTargetNode    = SettingTargetType("node")
	SettingTargetProject = SettingTargetType("project")
	SettingTargetApp     = SettingTargetType("app")
	SettingTargetEnvVar  = SettingTargetType("env-var")
)

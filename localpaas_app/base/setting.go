package base

type SettingTargetType string

const (
	SettingTargetProject = SettingTargetType("project")
	SettingTargetApp     = SettingTargetType("app")
	SettingTargetEnvVar  = SettingTargetType("env-var")
)

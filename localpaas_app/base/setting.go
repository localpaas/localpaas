package base

type SettingTargetType string

const (
	SettingTargetProject    = SettingTargetType("project")
	SettingTargetProjectEnv = SettingTargetType("project-env")
	SettingTargetApp        = SettingTargetType("app")
	SettingTargetEnvVar     = SettingTargetType("env-var")
)

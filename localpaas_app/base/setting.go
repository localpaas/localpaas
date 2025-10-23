package base

type SettingType string

const (
	SettingTypeNode      = SettingType("node")
	SettingTypeProject   = SettingType("project")
	SettingTypeApp       = SettingType("app")
	SettingTypeEnvVar    = SettingType("env-var")
	SettingTypeS3Storage = SettingType("s3-storage")
	SettingTypeSSHKey    = SettingType("ssh-key")
	SettingTypeAPIKey    = SettingType("api-key")
)

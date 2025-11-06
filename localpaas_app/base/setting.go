package base

type SettingType string

const (
	SettingTypeProject    = SettingType("project")
	SettingTypeApp        = SettingType("app")
	SettingTypeDeployment = SettingType("deployment")
	SettingTypeEnvVar     = SettingType("env-var")
	SettingTypeS3Storage  = SettingType("s3-storage")
	SettingTypeSSHKey     = SettingType("ssh-key")
	SettingTypeAPIKey     = SettingType("api-key")
)

var (
	AllAppSettingTypes = []SettingType{SettingTypeApp, SettingTypeDeployment, SettingTypeEnvVar}

	AllProjectSettingTypes = []SettingType{SettingTypeProject, SettingTypeEnvVar}
)

type SettingStatus string

const (
	SettingStatusActive   = SettingStatus("active")
	SettingStatusPending  = SettingStatus("pending")
	SettingStatusDisabled = SettingStatus("disabled")
)

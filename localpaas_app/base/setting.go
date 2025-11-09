package base

type SettingType string

const (
	SettingTypeProject    = SettingType("project")
	SettingTypeApp        = SettingType("app")
	SettingTypeDeployment = SettingType("deployment")
	SettingTypeEnvVar     = SettingType("env-var")
	SettingTypeSecret     = SettingType("secret")
	SettingTypeS3Storage  = SettingType("s3-storage")
	SettingTypeOAuth      = SettingType("oauth")
	SettingTypeSSHKey     = SettingType("ssh-key")
	SettingTypeAPIKey     = SettingType("api-key")
	SettingTypeSlack      = SettingType("slack")
)

var (
	AllAppSettingTypes = []SettingType{SettingTypeApp, SettingTypeDeployment, SettingTypeEnvVar, SettingTypeSecret}

	AllProjectSettingTypes = []SettingType{SettingTypeProject, SettingTypeEnvVar, SettingTypeSecret}
)

type SettingStatus string

const (
	SettingStatusActive   = SettingStatus("active")
	SettingStatusPending  = SettingStatus("pending")
	SettingStatusDisabled = SettingStatus("disabled")
)

var (
	AllSettingStatuses = []SettingStatus{SettingStatusActive, SettingStatusPending, SettingStatusDisabled}
)

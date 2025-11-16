package base

type SettingType string

const (
	SettingTypeProject      = SettingType("project")
	SettingTypeApp          = SettingType("app")
	SettingTypeServiceSpec  = SettingType("service-spec")
	SettingTypeDeployment   = SettingType("deployment")
	SettingTypeEnvVar       = SettingType("env-var")
	SettingTypeSecret       = SettingType("secret")
	SettingTypeS3Storage    = SettingType("s3-storage")
	SettingTypeOAuth        = SettingType("oauth")
	SettingTypeSSHKey       = SettingType("ssh-key")
	SettingTypeAPIKey       = SettingType("api-key")
	SettingTypeSlack        = SettingType("slack")
	SettingTypeRegistryAuth = SettingType("registry-auth")
	SettingTypeSsl          = SettingType("ssl")
)

var (
	AllAppSettingTypes = []SettingType{SettingTypeApp, SettingTypeServiceSpec, SettingTypeDeployment,
		SettingTypeEnvVar, SettingTypeSecret}

	AllProjectSettingTypes = []SettingType{SettingTypeProject, SettingTypeEnvVar, SettingTypeSecret}
)

type SettingStatus string

const (
	SettingStatusActive   = SettingStatus("active")
	SettingStatusPending  = SettingStatus("pending")
	SettingStatusDisabled = SettingStatus("disabled")
	SettingStatusExpired  = SettingStatus("expired")
)

var (
	AllSettingStatuses = []SettingStatus{SettingStatusActive, SettingStatusPending, SettingStatusDisabled}
)

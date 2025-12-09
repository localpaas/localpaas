package base

import "github.com/tiendc/gofn"

type SettingType string

const (
	SettingTypeProject       SettingType = "project"
	SettingTypeApp           SettingType = "app"
	SettingTypeServiceSpec   SettingType = "service-spec"
	SettingTypeAppDeployment SettingType = "app-deployment"
	SettingTypeAppHttp       SettingType = "app-http"
	SettingTypeEnvVar        SettingType = "env-var"
	SettingTypeSecret        SettingType = "secret"
	SettingTypeS3Storage     SettingType = "s3-storage"
	SettingTypeOAuth         SettingType = "oauth"
	SettingTypeSSHKey        SettingType = "ssh-key"
	SettingTypeAPIKey        SettingType = "api-key"
	SettingTypeSlack         SettingType = "slack"
	SettingTypeDiscord       SettingType = "discord"
	SettingTypeRegistryAuth  SettingType = "registry-auth"
	SettingTypeBasicAuth     SettingType = "basic-auth"
	SettingTypeSsl           SettingType = "ssl"
	SettingTypeGithubApp     SettingType = "github-app"
	SettingTypeGitToken      SettingType = "git-token"
	SettingTypeCronJob       SettingType = "cron-job"
)

var (
	AllAppSettingTypes = []SettingType{SettingTypeApp, SettingTypeAppDeployment,
		SettingTypeAppHttp, SettingTypeEnvVar, SettingTypeSecret}

	AllProjectSettingTypes = []SettingType{SettingTypeProject, SettingTypeEnvVar, SettingTypeSecret}
)

type SettingStatus string

const (
	SettingStatusActive   SettingStatus = "active"
	SettingStatusPending  SettingStatus = "pending"
	SettingStatusDisabled SettingStatus = "disabled"
	SettingStatusExpired  SettingStatus = "expired"
)

var (
	AllSettingStatuses = []SettingStatus{SettingStatusActive, SettingStatusPending, SettingStatusDisabled,
		SettingStatusExpired}
	AllSettingSettableStatuses = gofn.Drop(AllSettingStatuses, SettingStatusExpired)
)

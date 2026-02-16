package base

import "github.com/tiendc/gofn"

type SettingType string

const (
	SettingTypeProject         SettingType = "project"
	SettingTypeApp             SettingType = "app"
	SettingTypeAppDeployment   SettingType = "app-deployment"
	SettingTypeAppHttp         SettingType = "app-http"
	SettingTypeAppNotification SettingType = "app-notification"
	SettingTypeEnvVar          SettingType = "env-var"
	SettingTypeSecret          SettingType = "secret"
	SettingTypeAWS             SettingType = "aws"
	SettingTypeAWSS3           SettingType = "aws-s3"
	SettingTypeOAuth           SettingType = "oauth"
	SettingTypeSSHKey          SettingType = "ssh-key"
	SettingTypeAPIKey          SettingType = "api-key"
	SettingTypeIMService       SettingType = "im-service"
	SettingTypeRegistryAuth    SettingType = "registry-auth"
	SettingTypeBasicAuth       SettingType = "basic-auth"
	SettingTypeSSL             SettingType = "ssl"
	SettingTypeGithubApp       SettingType = "github-app"
	SettingTypeAccessToken     SettingType = "access-token"
	SettingTypeCronJob         SettingType = "cron-job"
	SettingTypeHealthcheck     SettingType = "healthcheck"
	SettingTypeEmail           SettingType = "email"
	SettingTypeRepoWebhook     SettingType = "repo-webhook"
	SettingTypeNotification    SettingType = "notification"
)

var (
	AllAppSettingTypes = []SettingType{SettingTypeApp, SettingTypeAppDeployment,
		SettingTypeAppHttp, SettingTypeEnvVar, SettingTypeSecret, SettingTypeCronJob, SettingTypeHealthcheck}

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

type SettingScope int8

const (
	SettingScopeGlobal SettingScope = iota
	SettingScopeUser
	SettingScopeProject
	SettingScopeApp
)

const (
	SettingNameMaxLen = 100
)

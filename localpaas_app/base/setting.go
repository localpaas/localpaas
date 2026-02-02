package base

import "github.com/tiendc/gofn"

type SettingType string

const (
	SettingTypeProject       SettingType = "project"
	SettingTypeApp           SettingType = "app"
	SettingTypeAppDeployment SettingType = "app-deployment"
	SettingTypeAppHttp       SettingType = "app-http"
	SettingTypeEnvVar        SettingType = "env-var"
	SettingTypeSecret        SettingType = "secret"
	SettingTypeAWS           SettingType = "aws"
	SettingTypeAWSS3         SettingType = "aws-s3"
	SettingTypeOAuth         SettingType = "oauth"
	SettingTypeSSHKey        SettingType = "ssh-key"
	SettingTypeAPIKey        SettingType = "api-key"
	SettingTypeIMService     SettingType = "im-service"
	SettingTypeRegistryAuth  SettingType = "registry-auth"
	SettingTypeBasicAuth     SettingType = "basic-auth"
	SettingTypeSSL           SettingType = "ssl"
	SettingTypeGithubApp     SettingType = "github-app"
	SettingTypeGitToken      SettingType = "git-token"
	SettingTypeCronJob       SettingType = "cron-job"
	SettingTypeEmail         SettingType = "email"
	SettingTypeWebhook       SettingType = "webhook"
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

type SettingScope int8

const (
	SettingScopeGlobal SettingScope = iota
	SettingScopeUser
	SettingScopeProject
	SettingScopeApp
)

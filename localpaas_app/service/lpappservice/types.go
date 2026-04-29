package lpappservice

import "github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"

type AppReleaseInfo struct {
	Stable *ReleaseInfo `json:"stable"`
	Beta   *ReleaseInfo `json:"beta"`
}

type ReleaseInfo struct {
	ReleaseDate  timeutil.Date `json:"releaseDate"`
	AppVersion   string        `json:"appVersion"`
	AppImage     string        `json:"appImage"`
	RedisImage   string        `json:"redisImage"`
	DbImage      string        `json:"dbImage"`
	TraefikImage string        `json:"traefikImage"`

	CanUpdate bool `json:"canUpdate"`
}

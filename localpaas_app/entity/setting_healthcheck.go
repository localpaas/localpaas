package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	CurrentHealthcheckVersion = 1
)

type Healthcheck struct {
	Type         base.HealthcheckType      `json:"type"`
	Interval     timeutil.Duration         `json:"interval"`
	MaxRetry     int                       `json:"maxRetry,omitempty"`
	RetryDelay   timeutil.Duration         `json:"retryDelay,omitempty"`
	Timeout      timeutil.Duration         `json:"timeout,omitempty"`
	REST         *HealthcheckREST          `json:"rest,omitempty"`
	GRPC         *HealthcheckGRPC          `json:"grpc,omitempty"`
	Notification *DefaultResultNtfnSetting `json:"notification,omitempty"`
}

type HealthcheckREST struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	ReturnCode  int    `json:"returnCode"`
	ReturnText  string `json:"returnText"`
	ReturnJSON  string `json:"returnJSON"`
}

type HealthcheckGRPC struct {
	// TODO: implement this
}

func (s *Healthcheck) GetType() base.SettingType {
	return base.SettingTypeHealthcheck
}

func (s *Healthcheck) GetRefSettingIDs() []string {
	res := make([]string, 0, 5) //nolint
	res = append(res, s.Notification.GetRefSettingIDs()...)
	return res
}

func (s *Setting) AsHealthcheck() (*Healthcheck, error) {
	return parseSettingAs(s, func() *Healthcheck { return &Healthcheck{} })
}

func (s *Setting) MustAsHealthcheck() *Healthcheck {
	return gofn.Must(s.AsHealthcheck())
}

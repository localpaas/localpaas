package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	CurrentHealthcheckVersion = 1
)

var _ = registerSettingParser(base.SettingTypeHealthcheck, &healthcheckParser{})

type healthcheckParser struct {
}

func (s *healthcheckParser) New() SettingData {
	return &Healthcheck{}
}

type Healthcheck struct {
	HealthcheckType base.HealthcheckType       `json:"healthcheckType"`
	Interval        timeutil.Duration          `json:"interval"`
	MaxRetry        int                        `json:"maxRetry,omitempty"`
	RetryDelay      timeutil.Duration          `json:"retryDelay,omitempty"`
	Timeout         timeutil.Duration          `json:"timeout,omitempty"`
	REST            *HealthcheckREST           `json:"rest,omitempty"`
	GRPC            *HealthcheckGRPC           `json:"grpc,omitempty"`
	Notification    *DefaultResultNotifSetting `json:"notification,omitempty"`
}

type HealthcheckREST struct {
	URL         string `json:"url"`
	Method      string `json:"method,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Body        string `json:"body,omitempty"`
	ReturnCode  int    `json:"returnCode,omitempty"`
	ReturnText  string `json:"returnText,omitempty"`
	ReturnJSON  string `json:"returnJSON,omitempty"`
}

type HealthcheckGRPC struct {
	Version      base.HealthcheckGRPCVersion `json:"version"`
	Addr         string                      `json:"addr"`
	Service      string                      `json:"service"`
	ReturnStatus base.HealthcheckGRPCStatus  `json:"returnStatus"`
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
	return parseSettingAs[*Healthcheck](s)
}

func (s *Setting) MustAsHealthcheck() *Healthcheck {
	return gofn.Must(s.AsHealthcheck())
}

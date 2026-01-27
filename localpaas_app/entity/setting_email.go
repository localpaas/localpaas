package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentEmailVersion = 1
)

type Email struct {
	SMTP *SMTPConf     `json:"smtp,omitempty"`
	HTTP *HTTPMailConf `json:"http,omitempty"`
}

type SMTPConf struct {
	Host        string         `json:"host"`
	Port        int            `json:"port"`
	Username    string         `json:"username"`
	DisplayName string         `json:"displayName"`
	Password    EncryptedField `json:"password"`
	SSL         bool           `json:"ssl"`
}

type HTTPMailConf struct {
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	BodyMapping map[string]string `json:"bodyMapping"`
}

func (s *Email) MustDecrypt() *Email {
	if s.SMTP != nil {
		s.SMTP.Password.MustGetPlain()
	}
	return s
}

func (s *Setting) AsEmail() (*Email, error) {
	return parseSettingAs(s, base.SettingTypeEmail, func() *Email { return &Email{} })
}

func (s *Setting) MustAsEmail() *Email {
	return gofn.Must(s.AsEmail())
}

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
	Endpoint     string                `json:"endpoint"`
	Method       string                `json:"method"`
	ContentType  string                `json:"contentType"`
	Headers      map[string]string     `json:"headers"`
	FieldMapping *HTTPMailFieldMapping `json:"fieldMapping"`
	Username     string                `json:"username"`
	DisplayName  string                `json:"displayName"`
	Password     EncryptedField        `json:"password"`
}

type HTTPMailFieldMapping struct {
	FromAddress string `json:"fromAddress"`
	FromName    string `json:"fromName"`
	ToAddress   string `json:"toAddress"`
	ToAddresses string `json:"toAddresses"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	Password    string `json:"password"`
}

func (s *Email) MustDecrypt() *Email {
	if s.SMTP != nil {
		s.SMTP.Password.MustGetPlain()
	}
	if s.HTTP != nil {
		s.HTTP.Password.MustGetPlain()
	}
	return s
}

func (s *Setting) AsEmail() (*Email, error) {
	return parseSettingAs(s, base.SettingTypeEmail, func() *Email { return &Email{} })
}

func (s *Setting) MustAsEmail() *Email {
	return gofn.Must(s.AsEmail())
}

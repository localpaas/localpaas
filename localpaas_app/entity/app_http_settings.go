package entity

import (
	"strings"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/services/nginx"
)

const (
	CurrentAppHttpSettingsVersion = 1
)

type AppHttpSettings struct {
	Enabled bool         `json:"enabled"`
	Domains []*AppDomain `json:"domains,omitempty"`
}

type AppDomain struct {
	Enabled          bool           `json:"enabled"`
	Domain           string         `json:"domain"`
	DomainRedirect   string         `json:"domainRedirect,omitempty"`
	SslCert          ObjectID       `json:"sslCert,omitzero"`
	ContainerPort    int            `json:"containerPort,omitempty"`
	ForceHttps       bool           `json:"forceHttps,omitempty"`
	WebsocketEnabled bool           `json:"websocketEnabled,omitempty"`
	BasicAuth        ObjectID       `json:"basicAuth,omitzero"`
	NginxSettings    *NginxSettings `json:"nginxSettings,omitempty"`
}

type NginxSettings struct {
	RootDirectives []*NginxDirective `json:"rootDirectives,omitempty"`
	ServerBlock    *NginxServerBlock `json:"serverBlock"`
}

type NginxServerBlock struct {
	Hide       bool              `json:"hide,omitempty"`
	Directives []*NginxDirective `json:"directives"`
}

type NginxDirective struct {
	Hide bool `json:"hide,omitempty"`
	*nginx.Directive
}

func (s *AppHttpSettings) GetDomain(domain string) *AppDomain {
	domain = strings.ToLower(domain)
	for _, domainRec := range s.Domains {
		if domainRec.Domain == domain {
			return domainRec
		}
	}
	return nil
}

func (s *AppHttpSettings) GetInUseSslCertIDs() (res []string) {
	for _, domain := range s.Domains {
		if !domain.Enabled {
			continue
		}
		if domain.SslCert.ID != "" {
			res = append(res, domain.SslCert.ID)
		}
	}
	return
}

func (s *AppHttpSettings) GetInUseBasicAuthIDs() (res []string) {
	for _, domain := range s.Domains {
		if !domain.Enabled {
			continue
		}
		if domain.BasicAuth.ID != "" {
			res = append(res, domain.BasicAuth.ID)
		}
	}
	return
}

func (s *AppHttpSettings) GetAllInUseSettingIDs() (res []string) {
	res = make([]string, 0, 5) //nolint
	res = append(res, s.GetInUseSslCertIDs()...)
	res = append(res, s.GetInUseBasicAuthIDs()...)
	return res
}

func (s *Setting) AsAppHttpSettings() (*AppHttpSettings, error) {
	return parseSettingAs(s, base.SettingTypeAppHttp, func() *AppHttpSettings { return &AppHttpSettings{} })
}

func (s *Setting) MustAsAppHttpSettings() *AppHttpSettings {
	return gofn.Must(s.AsAppHttpSettings())
}

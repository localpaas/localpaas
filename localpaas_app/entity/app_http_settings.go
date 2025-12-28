package entity

import (
	"strings"

	crossplane "github.com/nginxinc/nginx-go-crossplane"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
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
	*crossplane.Directive
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

func (s *Setting) AsAppHttpSettings() (*AppHttpSettings, error) {
	return parseSettingAs(s, base.SettingTypeAppHttp, func() *AppHttpSettings { return &AppHttpSettings{} })
}

func (s *Setting) MustAsAppHttpSettings() *AppHttpSettings {
	return gofn.Must(s.AsAppHttpSettings())
}

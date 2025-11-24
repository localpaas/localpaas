package entity

import (
	"strings"

	crossplane "github.com/nginxinc/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type AppHttpSettings struct {
	Enabled bool         `json:"enabled"`
	Domains []*AppDomain `json:"domains,omitempty"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
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

func (s *Setting) ParseAppHttpSettings() (*AppHttpSettings, error) {
	res := &AppHttpSettings{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeAppHttp {
		return res, s.parseData(res)
	}
	return nil, nil
}

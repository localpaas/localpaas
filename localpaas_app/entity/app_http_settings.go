package entity

import (
	"strings"

	crossplane "github.com/nginxinc/nginx-go-crossplane"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type AppHttpSettings struct {
	Enabled          bool           `json:"enabled"`
	Domains          []*AppDomain   `json:"domains,omitempty"`
	DomainRedirect   string         `json:"domainRedirect,omitempty"`
	ContainerPort    int            `json:"containerPort,omitempty"`
	ForceHttps       bool           `json:"forceHttps,omitempty"`
	WebsocketEnabled bool           `json:"websocketEnabled,omitempty"`
	BasicAuth        ObjectID       `json:"basicAuth,omitzero"`
	NginxSettings    *NginxSettings `json:"nginxSettings,omitempty"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

type AppDomain struct {
	Domain  string   `json:"domain"`
	SslCert ObjectID `json:"sslCert,omitzero"`
}

type NginxSettings struct {
	Enabled        bool              `json:"enabled"`
	RootDirectives []*NginxDirective `json:"rootDirectives,omitempty"`
	ServerBlock    *NginxServerBlock `json:"serverBlock"`
}

type NginxServerBlock struct {
	Invisible  bool              `json:"invisible,omitempty"`
	Directives []*NginxDirective `json:"directives"`
}

type NginxDirective struct {
	Invisible bool `json:"invisible,omitempty"`
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

package appdto

import (
	crossplane "github.com/nginxinc/nginx-go-crossplane"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

//
// REQUEST
//

type HttpSettingsReq struct {
	Enabled          bool                `json:"enabled"`
	Domains          []*DomainReq        `json:"domains"`
	DomainRedirect   string              `json:"domainRedirect"`
	ContainerPort    int                 `json:"containerPort"`
	ForceHttps       bool                `json:"forceHttps"`
	WebsocketEnabled bool                `json:"websocketEnabled"`
	BasicAuth        basedto.ObjectIDReq `json:"basicAuth"`
	NginxSettings    *NginxSettingsReq   `json:"nginxSettings"`
}

// nolint
func (req *HttpSettingsReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	// TODO:
	return res
}

type DomainReq struct {
	Domain  string              `json:"domain"`
	SslCert basedto.ObjectIDReq `json:"sslCert"`
}

// nolint
func (req *DomainReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return
	}
	if field != "" {
		field += "."
	}
	// TODO:
	return res
}

type NginxSettingsReq struct {
	Enabled        bool                 `json:"enabled"`
	RootDirectives []*NginxDirectiveReq `json:"rootDirectives"`
	ServerBlock    *NginxServerBlockReq `json:"serverBlock"`
}

type NginxServerBlockReq struct {
	Invisible  bool                 `json:"invisible"`
	Directives []*NginxDirectiveReq `json:"directives"`
}

type NginxDirectiveReq struct {
	Invisible bool `json:"invisible"`
	*crossplane.Directive
}

//
// RESPONSE
//

type HttpSettingsResp struct {
	Enabled          bool                     `json:"enabled"`
	Domains          []*DomainResp            `json:"domains"`
	DomainRedirect   string                   `json:"domainRedirect"`
	ContainerPort    int                      `json:"containerPort"`
	ForceHttps       bool                     `json:"forceHttps"`
	WebsocketEnabled bool                     `json:"websocketEnabled"`
	BasicAuth        *basedto.NamedObjectResp `json:"basicAuth"`
	NginxSettings    *NginxSettingsResp       `json:"nginxSettings"`
}

type DomainResp struct {
	Domain  string                   `json:"domain"`
	SslCert *basedto.NamedObjectResp `json:"sslCert"`
}

type NginxSettingsResp struct {
	Enabled        bool                  `json:"enabled"`
	RootDirectives []*NginxDirectiveResp `json:"rootDirectives"`
	ServerBlock    *NginxServerBlockResp `json:"serverBlock"`
}

type NginxServerBlockResp struct {
	Invisible  bool                  `json:"invisible"`
	Directives []*NginxDirectiveResp `json:"directives"`
}

type NginxDirectiveResp struct {
	Invisible bool `json:"invisible"`
	*crossplane.Directive
}

func TransformHttpSettings(setting *entity.Setting) (resp *HttpSettingsResp, err error) {
	data, err := setting.ParseAppHttpSettings()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, &data); err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

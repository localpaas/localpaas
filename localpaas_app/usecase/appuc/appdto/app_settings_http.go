package appdto

import (
	crossplane "github.com/nginxinc/nginx-go-crossplane"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

//
// REQUEST
//

type HttpSettingsReq struct {
	Enabled bool         `json:"enabled"`
	Domains []*DomainReq `json:"domains"`
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
	Enabled          bool                `json:"enabled"`
	Domain           string              `json:"domain"`
	DomainRedirect   string              `json:"domainRedirect"`
	SslCert          basedto.ObjectIDReq `json:"sslCert"`
	ContainerPort    int                 `json:"containerPort"`
	ForceHttps       bool                `json:"forceHttps"`
	WebsocketEnabled bool                `json:"websocketEnabled"`
	BasicAuth        basedto.ObjectIDReq `json:"basicAuth"`
	NginxSettings    *NginxSettingsReq   `json:"nginxSettings"`
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
	RootDirectives []*NginxDirectiveReq `json:"rootDirectives"`
	ServerBlock    *NginxServerBlockReq `json:"serverBlock"`
}

type NginxServerBlockReq struct {
	Hide       bool                 `json:"hide"`
	Directives []*NginxDirectiveReq `json:"directives"`
}

type NginxDirectiveReq struct {
	Hide bool `json:"hide"`
	*crossplane.Directive
}

//
// RESPONSE
//

type HttpSettingsResp struct {
	Enabled              bool               `json:"enabled"`
	Domains              []*DomainResp      `json:"domains"`
	DefaultNginxSettings *NginxSettingsResp `json:"defaultNginxSettings"`
}

type DomainResp struct {
	Enabled          bool                        `json:"enabled"`
	Domain           string                      `json:"domain"`
	DomainRedirect   string                      `json:"domainRedirect"`
	SslCert          *ssldto.SslResp             `json:"sslCert"`
	ContainerPort    int                         `json:"containerPort"`
	ForceHttps       bool                        `json:"forceHttps"`
	WebsocketEnabled bool                        `json:"websocketEnabled"`
	BasicAuth        *basicauthdto.BasicAuthResp `json:"basicAuth"`
	NginxSettings    *NginxSettingsResp          `json:"nginxSettings"`
}

type NginxSettingsResp struct {
	RootDirectives []*NginxDirectiveResp `json:"rootDirectives,omitempty"`
	ServerBlock    *NginxServerBlockResp `json:"serverBlock"`
}

type NginxServerBlockResp struct {
	Hide       bool                  `json:"hide,omitempty"`
	Directives []*NginxDirectiveResp `json:"directives"`
}

type NginxDirectiveResp struct {
	Hide bool `json:"hide,omitempty"`
	*crossplane.Directive
}

func TransformHttpSettings(input *AppSettingsTransformationInput) (resp *HttpSettingsResp, err error) {
	if input.HttpSettings == nil {
		return nil, nil
	}

	if err = copier.Copy(&resp, input.HttpSettings); err != nil {
		return nil, apperrors.Wrap(err)
	}

	for _, domain := range resp.Domains {
		if domain.SslCert != nil && domain.SslCert.ID != "" {
			sslResp, _ := ssldto.TransformSsl(input.ReferenceSettingMap[domain.SslCert.ID])
			if sslResp != nil {
				domain.SslCert = sslResp
			}
		} else {
			domain.SslCert = nil
		}

		if domain.BasicAuth != nil && domain.BasicAuth.ID != "" {
			basicAuthResp, _ := basicauthdto.TransformBasicAuth(input.ReferenceSettingMap[domain.BasicAuth.ID])
			if basicAuthResp != nil {
				domain.BasicAuth = basicAuthResp
			}
		} else {
			domain.BasicAuth = nil
		}
	}

	if err = copier.Copy(&resp.DefaultNginxSettings, input.DefaultNginxSettings); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
	"github.com/localpaas/localpaas/services/nginx"
)

type GetAppHttpSettingsReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
}

func NewGetAppHttpSettingsReq() *GetAppHttpSettingsReq {
	return &GetAppHttpSettingsReq{}
}

func (req *GetAppHttpSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAppHttpSettingsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *HttpSettingsResp `json:"data"`
}

type HttpSettingsResp struct {
	Enabled              bool               `json:"enabled"`
	Domains              []*DomainResp      `json:"domains"`
	DefaultNginxSettings *NginxSettingsResp `json:"defaultNginxSettings"`
	UpdateVer            int                `json:"updateVer"`
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
	*nginx.Directive
}

type AppHttpSettingsTransformInput struct {
	App          *entity.App
	HttpSettings *entity.Setting

	DefaultNginxSettings *entity.NginxSettings
	ReferenceSettingMap  map[string]*entity.Setting
}

func TransformHttpSettings(input *AppHttpSettingsTransformInput) (resp *HttpSettingsResp, err error) {
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

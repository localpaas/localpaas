package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
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
	Enabled          bool                      `json:"enabled"`
	Domain           string                    `json:"domain"`
	DomainRedirect   string                    `json:"domainRedirect"`
	SslCert          *settings.BaseSettingResp `json:"sslCert"`
	ContainerPort    int                       `json:"containerPort"`
	ForceHttps       bool                      `json:"forceHttps"`
	WebsocketEnabled bool                      `json:"websocketEnabled"`
	BasicAuth        *settings.BaseSettingResp `json:"basicAuth"`
	NginxSettings    *NginxSettingsResp        `json:"nginxSettings"`
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
	RefSettingMap        map[string]*entity.Setting
}

func TransformHttpSettings(input *AppHttpSettingsTransformInput) (resp *HttpSettingsResp, err error) {
	if input.HttpSettings == nil {
		return nil, nil
	}

	if err = copier.Copy(&resp, input.HttpSettings); err != nil {
		return nil, apperrors.Wrap(err)
	}
	appHttpSettings := input.HttpSettings.MustAsAppHttpSettings()
	if err = copier.Copy(&resp, appHttpSettings); err != nil {
		return nil, apperrors.Wrap(err)
	}

	for _, domain := range resp.Domains {
		if domain.SslCert != nil && domain.SslCert.ID != "" {
			settingResp, _ := settings.TransformSettingBase(input.RefSettingMap[domain.SslCert.ID])
			if settingResp != nil {
				domain.SslCert = settingResp
			} else {
				domain.SslCert = nil
			}
		}

		if domain.BasicAuth != nil && domain.BasicAuth.ID != "" {
			settingResp, _ := settings.TransformSettingBase(input.RefSettingMap[domain.BasicAuth.ID])
			if settingResp != nil {
				domain.BasicAuth = settingResp
			} else {
				domain.BasicAuth = nil
			}
		}
	}

	if err = copier.Copy(&resp.DefaultNginxSettings, input.DefaultNginxSettings); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

package appdto

import (
	crossplane "github.com/nginxinc/nginx-go-crossplane"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateAppHttpSettingsReq struct {
	ProjectID string       `json:"-"`
	AppID     string       `json:"-"`
	Enabled   bool         `json:"enabled"`
	Domains   []*DomainReq `json:"domains"`
	UpdateVer int          `json:"updateVer"`
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

func NewUpdateAppHttpSettingsReq() *UpdateAppHttpSettingsReq {
	return &UpdateAppHttpSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateAppHttpSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	// TODO: validate http settings input
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateAppHttpSettingsResp struct {
	Meta *basedto.BaseMeta              `json:"meta"`
	Data *UpdateAppHttpSettingsDataResp `json:"data"`
}

type UpdateAppHttpSettingsDataResp struct {
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

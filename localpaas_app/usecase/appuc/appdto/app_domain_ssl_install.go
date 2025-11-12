package appdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	domainMaxLen = 100
)

type InstallDomainSslReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
	Domain    string `json:"domain"`
}

func NewInstallDomainSslReq() *InstallDomainSslReq {
	return &InstallDomainSslReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *InstallDomainSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, basedto.ValidateStr(&req.Domain, true, 1, domainMaxLen, "domain")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type InstallDomainSslResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

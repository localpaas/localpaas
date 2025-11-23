package appdto

import (
	"strings"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	domainMaxLen = 100
)

type ObtainDomainSslReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
	Domain    string `json:"domain"`
}

func NewObtainDomainSslReq() *ObtainDomainSslReq {
	return &ObtainDomainSslReq{}
}

func (req *ObtainDomainSslReq) ModifyRequest() error {
	req.Domain = strings.TrimSpace(strings.ToLower(req.Domain))
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *ObtainDomainSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, basedto.ValidateStr(&req.Domain, true, 1, domainMaxLen, "domain")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ObtainDomainSslResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

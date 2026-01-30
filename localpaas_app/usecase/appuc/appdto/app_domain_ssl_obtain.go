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

var (
	allowedKeySizes = []int{0, 2048, 3072, 4096}
)

type ObtainDomainSSLReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`
	Domain    string `json:"domain"`
	Email     string `json:"email"`
	KeySize   int    `json:"keySize"`
}

func NewObtainDomainSSLReq() *ObtainDomainSSLReq {
	return &ObtainDomainSSLReq{}
}

func (req *ObtainDomainSSLReq) ModifyRequest() error {
	req.Domain = strings.TrimSpace(strings.ToLower(req.Domain))
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *ObtainDomainSSLReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, basedto.ValidateStr(&req.Domain, true, 1, domainMaxLen, "domain")...)
	validators = append(validators, basedto.ValidateEmail(&req.Email, false, "email")...)
	validators = append(validators, basedto.ValidateNumberIn(&req.KeySize, false, allowedKeySizes, "keySize")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ObtainDomainSSLResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

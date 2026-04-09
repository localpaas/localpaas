package sslcertsettingsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateUniqueSSLCertSettingsReq struct {
	settings.UpdateUniqueSettingReq
	*SSLCertSettingsBaseReq
}

type SSLCertSettingsBaseReq struct {
	CertType    base.SSLCertType  `json:"certType"`
	KeyType     base.SSLKeyType   `json:"keyType"`
	ValidPeriod timeutil.Duration `json:"validPeriod"`
	RootDomain  string            `json:"rootDomain"`
	Email       string            `json:"email"`
	AutoRenew   bool              `json:"autoRenew"`
}

func (req *SSLCertSettingsBaseReq) ToEntity() *entity.SSLCertSettings {
	return &entity.SSLCertSettings{
		CertType:    req.CertType,
		KeyType:     req.KeyType,
		ValidPeriod: req.ValidPeriod,
		RootDomain:  req.RootDomain,
		Email:       req.Email,
		AutoRenew:   req.AutoRenew,
	}
}

// nolint
func (req *SSLCertSettingsBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	// TODO: add validation
	return res
}

func NewUpdateUniqueSSLCertSettingsReq() *UpdateUniqueSSLCertSettingsReq {
	return &UpdateUniqueSSLCertSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateUniqueSSLCertSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateUniqueSSLCertSettingsResp struct {
	Meta *basedto.Meta `json:"meta"`
}

package sslcertsettingsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteUniqueSSLCertSettingsReq struct {
	settings.DeleteUniqueSettingReq
}

func NewDeleteUniqueSSLCertSettingsReq() *DeleteUniqueSSLCertSettingsReq {
	return &DeleteUniqueSSLCertSettingsReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteUniqueSSLCertSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteUniqueSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteUniqueSSLCertSettingsResp struct {
	Meta *basedto.Meta `json:"meta"`
}

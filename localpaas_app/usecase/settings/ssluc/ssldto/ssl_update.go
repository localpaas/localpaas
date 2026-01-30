package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateSSLReq struct {
	settings.UpdateSettingReq
	*SSLBaseReq
}

func NewUpdateSSLReq() *UpdateSSLReq {
	return &UpdateSSLReq{}
}

func (req *UpdateSSLReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSSLReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSSLResp struct {
	Meta *basedto.Meta `json:"meta"`
}

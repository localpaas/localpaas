package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DeleteSslReq struct {
	settings.DeleteSettingReq
}

func NewDeleteSslReq() *DeleteSslReq {
	return &DeleteSslReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *DeleteSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.DeleteSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DeleteSslResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

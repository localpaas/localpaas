package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateSslReq struct {
	ID string `json:"-"`
	*SslBaseReq
}

func NewUpdateSslReq() *UpdateSslReq {
	return &UpdateSslReq{}
}

func (req *UpdateSslReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSslResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

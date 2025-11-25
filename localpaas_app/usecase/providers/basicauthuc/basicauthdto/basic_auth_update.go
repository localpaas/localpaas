package basicauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type UpdateBasicAuthReq struct {
	ID        string `json:"-"`
	UpdateVer int    `json:"updateVer"`
	*BasicAuthBaseReq
}

func NewUpdateBasicAuthReq() *UpdateBasicAuthReq {
	return &UpdateBasicAuthReq{}
}

func (req *UpdateBasicAuthReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateBasicAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateBasicAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

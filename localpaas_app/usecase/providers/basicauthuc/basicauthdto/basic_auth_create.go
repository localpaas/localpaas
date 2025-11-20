package basicauthdto

import (
	"strings"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	nameMaxLen = 100
)

type CreateBasicAuthReq struct {
	*BasicAuthBaseReq
}

type BasicAuthBaseReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *BasicAuthBaseReq) modifyRequest() error {
	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	return nil
}

func (req *BasicAuthBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, nameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.Username, true, 1, nameMaxLen, field+"username")...)
	res = append(res, basedto.ValidateStr(&req.Password, true, 1, nameMaxLen, field+"password")...)
	return res
}

func NewCreateBasicAuthReq() *CreateBasicAuthReq {
	return &CreateBasicAuthReq{}
}

func (req *CreateBasicAuthReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *CreateBasicAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateBasicAuthResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

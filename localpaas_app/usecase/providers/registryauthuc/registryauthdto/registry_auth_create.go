package registryauthdto

import (
	"strings"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

const (
	nameMaxLen = 100
)

type CreateRegistryAuthReq struct {
	providers.CreateSettingReq
	*RegistryAuthBaseReq
}

type RegistryAuthBaseReq struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *RegistryAuthBaseReq) modifyRequest() error {
	req.Name = strings.TrimSpace(req.Name)
	req.Address = strings.ToLower(strings.TrimSpace(req.Address))
	req.Username = strings.TrimSpace(req.Username)
	return nil
}

func (req *RegistryAuthBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Name, true, 1, nameMaxLen, field+"name")...)
	res = append(res, basedto.ValidateStr(&req.Address, true, 1, nameMaxLen, field+"address")...)
	res = append(res, basedto.ValidateStr(&req.Username, true, 1, nameMaxLen, field+"username")...)
	res = append(res, basedto.ValidateStr(&req.Password, true, 1, nameMaxLen, field+"password")...)
	return res
}

func NewCreateRegistryAuthReq() *CreateRegistryAuthReq {
	return &CreateRegistryAuthReq{}
}

func (req *CreateRegistryAuthReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *CreateRegistryAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateRegistryAuthResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

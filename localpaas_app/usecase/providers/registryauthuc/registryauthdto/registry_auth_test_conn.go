package registryauthdto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestRegistryAuthConnReq struct {
	*RegistryAuthBaseReq
}

func NewTestRegistryAuthConnReq() *TestRegistryAuthConnReq {
	return &TestRegistryAuthConnReq{}
}

func (req *TestRegistryAuthConnReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *TestRegistryAuthConnReq) Validate() apperrors.ValidationErrors {
	// NOTE: make sure req.Name is not empty to not fail the validation
	req.Name = gofn.Coalesce(req.Name, "x")

	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestRegistryAuthConnResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

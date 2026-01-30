package awss3dto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestAWSS3ConnReq struct {
	*AWSS3BaseReq
}

func NewTestAWSS3ConnReq() *TestAWSS3ConnReq {
	return &TestAWSS3ConnReq{}
}

func (req *TestAWSS3ConnReq) ModifyRequest() error {
	// NOTE: make sure req.Name is not empty to not fail the validation
	req.Name = gofn.Coalesce(req.Name, "x")
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *TestAWSS3ConnReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestAWSS3ConnResp struct {
	Meta *basedto.Meta `json:"meta"`
}

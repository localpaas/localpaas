package gittokendto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestGitTokenConnReq struct {
	*GitTokenBaseReq
}

func NewTestGitTokenConnReq() *TestGitTokenConnReq {
	return &TestGitTokenConnReq{}
}

func (req *TestGitTokenConnReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *TestGitTokenConnReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestGitTokenConnResp struct {
	Meta *basedto.Meta `json:"meta"`
}

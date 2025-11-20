package slackdto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestSendSlackMsgReq struct {
	*SlackBaseReq
	TestMsg string `json:"testMsg"`
}

func NewTestSendSlackMsgReq() *TestSendSlackMsgReq {
	return &TestSendSlackMsgReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *TestSendSlackMsgReq) Validate() apperrors.ValidationErrors {
	// NOTE: make sure req.Name is not empty to not fail the validation
	req.Name = gofn.Coalesce(req.Name, "x")
	req.TestMsg = gofn.Coalesce(req.TestMsg, "test message")

	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestSendSlackMsgResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

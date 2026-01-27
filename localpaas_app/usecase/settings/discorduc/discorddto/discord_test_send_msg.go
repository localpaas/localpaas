package discorddto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestSendDiscordMsgReq struct {
	*DiscordBaseReq
	TestMsg string `json:"testMsg"`
}

func NewTestSendDiscordMsgReq() *TestSendDiscordMsgReq {
	return &TestSendDiscordMsgReq{}
}

func (req *TestSendDiscordMsgReq) ModifyRequest() error {
	// NOTE: make sure req.Name is not empty to not fail the validation
	req.Name = gofn.Coalesce(req.Name, "x")
	req.TestMsg = gofn.Coalesce(req.TestMsg, "test message")
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *TestSendDiscordMsgReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestSendDiscordMsgResp struct {
	Meta *basedto.Meta `json:"meta"`
}

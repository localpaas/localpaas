package slackdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateSlackReq struct {
	settings.CreateSettingReq
	*SlackBaseReq
}

type SlackBaseReq struct {
	Name    string `json:"name"`
	Webhook string `json:"webhook"`
}

func (req *SlackBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateSlackReq() *CreateSlackReq {
	return &CreateSlackReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateSlackReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateSlackResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

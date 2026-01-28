package imservicedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateIMServiceReq struct {
	settings.CreateSettingReq
	*IMServiceBaseReq
}

type IMServiceBaseReq struct {
	Name    string      `json:"name"`
	Slack   *SlackReq   `json:"slack"`
	Discord *DiscordReq `json:"discord"`
}

type SlackReq struct {
	Webhook string `json:"webhook"`
}

type DiscordReq struct {
	Webhook string `json:"webhook"`
}

func (req *IMServiceBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateIMServiceReq() *CreateIMServiceReq {
	return &CreateIMServiceReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateIMServiceReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateIMServiceResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

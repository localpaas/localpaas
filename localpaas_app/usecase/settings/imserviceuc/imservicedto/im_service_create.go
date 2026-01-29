package imservicedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateIMServiceReq struct {
	settings.CreateSettingReq
	*IMServiceBaseReq
}

type IMServiceBaseReq struct {
	Name    string             `json:"name"`
	Kind    base.IMServiceKind `json:"kind"`
	Slack   *SlackReq          `json:"slack"`
	Discord *DiscordReq        `json:"discord"`
}

func (req *IMServiceBaseReq) ToEntity() *entity.IMService {
	imService := &entity.IMService{}
	switch req.Kind {
	case base.IMServiceKindSlack:
		imService.Slack = req.Slack.ToEntity()
	case base.IMServiceKindDiscord:
		imService.Discord = req.Discord.ToEntity()
	}
	return imService
}

type SlackReq struct {
	Webhook string `json:"webhook"`
}

func (req *SlackReq) ToEntity() *entity.Slack {
	return &entity.Slack{
		Webhook: entity.NewEncryptedField(req.Webhook),
	}
}

type DiscordReq struct {
	Webhook string `json:"webhook"`
}

func (req *DiscordReq) ToEntity() *entity.Discord {
	return &entity.Discord{
		Webhook: entity.NewEncryptedField(req.Webhook),
	}
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

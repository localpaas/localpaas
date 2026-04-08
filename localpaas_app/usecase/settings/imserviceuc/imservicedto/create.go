package imservicedto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	webhookURLMaxLen = 512
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

func (req *SlackReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return nil
	}
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Webhook, true, 1, webhookURLMaxLen, field+"webhook")...)
	return res
}

type DiscordReq struct {
	Webhook string `json:"webhook"`
}

func (req *DiscordReq) ToEntity() *entity.Discord {
	return &entity.Discord{
		Webhook: entity.NewEncryptedField(req.Webhook),
	}
}

func (req *DiscordReq) validate(field string) (res []vld.Validator) {
	if req == nil {
		return nil
	}
	if field != "" {
		field += "."
	}
	res = append(res, basedto.ValidateStr(&req.Webhook, true, 1, webhookURLMaxLen, field+"webhook")...)
	return res
}

func (req *IMServiceBaseReq) validate(field string) (res []vld.Validator) {
	if field != "" {
		field += "."
	}
	switch req.Kind {
	case base.IMServiceKindSlack:
		res = append(res, basedto.ValidateValue(req.Slack != nil, field+"slack")...)
		res = append(res, req.Slack.validate(field+"slack")...)
	case base.IMServiceKindDiscord:
		res = append(res, basedto.ValidateValue(req.Discord != nil, field+"discord")...)
		res = append(res, req.Discord.validate(field+"discord")...)
	}
	return res
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

package discorddto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateDiscordReq struct {
	*DiscordBaseReq
}

type DiscordBaseReq struct {
	Name    string `json:"name"`
	Webhook string `json:"webhook"`
}

func (req *DiscordBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateDiscordReq() *CreateDiscordReq {
	return &CreateDiscordReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateDiscordReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateDiscordResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

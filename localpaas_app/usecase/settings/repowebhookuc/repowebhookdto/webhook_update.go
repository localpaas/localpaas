package repowebhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateRepoWebhookReq struct {
	settings.UpdateSettingReq
	*RepoWebhookBaseReq
}

func NewUpdateRepoWebhookReq() *UpdateRepoWebhookReq {
	return &UpdateRepoWebhookReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateRepoWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateRepoWebhookResp struct {
	Meta *basedto.Meta `json:"meta"`
}

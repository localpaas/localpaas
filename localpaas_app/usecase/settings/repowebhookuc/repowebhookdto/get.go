package repowebhookdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetRepoWebhookReq struct {
	settings.GetSettingReq
}

func NewGetRepoWebhookReq() *GetRepoWebhookReq {
	return &GetRepoWebhookReq{}
}

func (req *GetRepoWebhookReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetRepoWebhookResp struct {
	Meta *basedto.Meta    `json:"meta"`
	Data *RepoWebhookResp `json:"data"`
}

type RepoWebhookResp struct {
	*settings.BaseSettingResp
	Kind   base.WebhookKind `json:"kind"`
	Secret string           `json:"secret"`
}

func TransformRepoWebhook(
	setting *entity.Setting,
	_ *entity.RefObjects,
) (resp *RepoWebhookResp, err error) {
	config := setting.MustAsRepoWebhook()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Kind = base.WebhookKind(setting.Kind)

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

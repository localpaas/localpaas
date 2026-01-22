package apikeydto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetAPIKeyReq struct {
	ID string `json:"-"`
}

func NewGetAPIKeyReq() *GetAPIKeyReq {
	return &GetAPIKeyReq{}
}

func (req *GetAPIKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAPIKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *APIKeyResp       `json:"data"`
}

type APIKeyResp struct {
	*settings.BaseSettingResp
	KeyID        string             `json:"keyId"`
	AccessAction base.AccessActions `json:"accessAction"`
}

func TransformAPIKey(setting *entity.Setting, objectID string) (resp *APIKeyResp, err error) {
	apiKey := setting.MustAsAPIKey()
	if err = copier.Copy(&resp, apiKey); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting, objectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

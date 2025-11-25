package apikeydto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
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
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Status       base.SettingStatus `json:"status"`
	KeyID        string             `json:"keyId"`
	AccessAction base.AccessActions `json:"accessAction"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformAPIKey(setting *entity.Setting) (resp *APIKeyResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	apiKey := setting.MustAsAPIKey()
	if err = copier.Copy(&resp, apiKey); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

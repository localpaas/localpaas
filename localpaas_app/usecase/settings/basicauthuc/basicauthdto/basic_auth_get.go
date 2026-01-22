package basicauthdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedPassword = "********"
)

type GetBasicAuthReq struct {
	settings.GetSettingReq
}

func NewGetBasicAuthReq() *GetBasicAuthReq {
	return &GetBasicAuthReq{}
}

func (req *GetBasicAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetBasicAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *BasicAuthResp    `json:"data"`
}

type BasicAuthResp struct {
	*settings.BaseSettingResp
	Username  string `json:"username"`
	Password  string `json:"password"`
	Encrypted bool   `json:"encrypted,omitempty"`
}

func (resp *BasicAuthResp) CopyPassword(field entity.EncryptedField) error {
	resp.Password = field.String()
	return nil
}

func TransformBasicAuth(setting *entity.Setting, objectID string) (resp *BasicAuthResp, err error) {
	config := setting.MustAsBasicAuth()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Encrypted = config.Password.IsEncrypted()
	if resp.Encrypted {
		resp.Password = maskedPassword
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting, objectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

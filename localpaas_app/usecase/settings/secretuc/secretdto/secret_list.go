package secretdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	maskedSecretValue = "****************"
)

type ListSecretReq struct {
	ObjectID string               `json:"-"`
	Status   []base.SettingStatus `json:"-" mapstructure:"status"`
	Search   string               `json:"-" mapstructure:"search"`

	Paging basedto.Paging `json:"-"`
}

func NewListSecretReq() *ListSecretReq {
	return &ListSecretReq{
		Paging: basedto.Paging{
			// Default paging if unset by client
			Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
		},
	}
}

func (req *ListSecretReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateSlice(req.Status, true, 0, base.AllSettingStatuses, "status")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListSecretResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*SecretResp `json:"data"`
}

type SecretResp struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Encrypted bool   `json:"encrypted,omitempty"`
}

func TransformSecret(setting *entity.Setting, decrypt bool) (resp *SecretResp, err error) {
	secret, err := setting.ParseSecret(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp = &SecretResp{
		ID:  setting.ID,
		Key: setting.Name,
	}
	if secret != nil {
		resp.Value = secret.Value
		resp.Encrypted = secret.IsEncrypted()
		if resp.Encrypted {
			resp.Value = maskedSecretValue
		}
	}
	return resp, nil
}

func TransformSecrets(settings []*entity.Setting, decrypt bool) (resp []*SecretResp, err error) {
	resp = make([]*SecretResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformSecret(setting, decrypt)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}

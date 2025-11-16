package secretdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
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

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformSecret(setting *entity.Setting, decrypt bool) (resp *SecretResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Key = setting.Name

	secret, err := setting.ParseSecret(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, &secret); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = secret.IsEncrypted()
	if resp.Encrypted {
		resp.Value = maskedSecretValue
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

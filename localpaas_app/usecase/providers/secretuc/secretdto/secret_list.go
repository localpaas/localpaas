package secretdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
)

type ListSecretReq struct {
	providers.ListSettingReq
}

func NewListSecretReq() *ListSecretReq {
	return &ListSecretReq{
		ListSettingReq: providers.ListSettingReq{
			Paging: basedto.Paging{
				// Default paging if unset by client
				Sort: basedto.Orders{{Direction: basedto.DirectionAsc, ColumnName: "name"}},
			},
		},
	}
}

func (req *ListSecretReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.ListSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListSecretResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*SecretResp `json:"data"`
}

type SecretResp struct {
	ID        string             `json:"id"`
	Name      string             `json:"name,omitempty"`
	Status    base.SettingStatus `json:"status"`
	Key       string             `json:"key"`
	UpdateVer int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformSecret(setting *entity.Setting) (resp *SecretResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Key = setting.Name

	secret := setting.MustAsSecret()
	if err = copier.Copy(&resp, &secret); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

func TransformSecrets(settings []*entity.Setting) (resp []*SecretResp, err error) {
	resp = make([]*SecretResp, 0, len(settings))
	for _, setting := range settings {
		item, err := TransformSecret(setting)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, item)
	}
	return resp, nil
}

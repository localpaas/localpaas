package registryauthdto

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
	maskedPassword = "********"
)

type GetRegistryAuthReq struct {
	ID string `json:"-"`
}

func NewGetRegistryAuthReq() *GetRegistryAuthReq {
	return &GetRegistryAuthReq{}
}

func (req *GetRegistryAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetRegistryAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *RegistryAuthResp `json:"data"`
}

type RegistryAuthResp struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Status    base.SettingStatus `json:"status"`
	Address   string             `json:"address"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Encrypted bool               `json:"encrypted,omitempty"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformRegistryAuth(setting *entity.Setting, decrypt bool) (resp *RegistryAuthResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Address = setting.Kind

	config, err := setting.ParseRegistryAuth(decrypt)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.IsEncrypted()
	if resp.Encrypted {
		resp.Password = maskedPassword
	}
	return resp, nil
}

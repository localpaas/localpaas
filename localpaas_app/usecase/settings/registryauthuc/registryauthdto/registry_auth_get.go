package registryauthdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedPassword = "********"
)

type GetRegistryAuthReq struct {
	settings.GetSettingReq
}

func NewGetRegistryAuthReq() *GetRegistryAuthReq {
	return &GetRegistryAuthReq{}
}

func (req *GetRegistryAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
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
	UpdateVer int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *RegistryAuthResp) CopyPassword(field entity.EncryptedField) error {
	resp.Password = field.String()
	return nil
}

func TransformRegistryAuth(setting *entity.Setting) (resp *RegistryAuthResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.Address = setting.Kind

	config := setting.MustAsRegistryAuth()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Password.IsEncrypted()
	if resp.Encrypted {
		resp.Password = maskedPassword
	}
	return resp, nil
}

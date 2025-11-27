package basicauthdto

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

type GetBasicAuthReq struct {
	ID string `json:"-"`
}

func NewGetBasicAuthReq() *GetBasicAuthReq {
	return &GetBasicAuthReq{}
}

func (req *GetBasicAuthReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetBasicAuthResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *BasicAuthResp    `json:"data"`
}

type BasicAuthResp struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Status    base.SettingStatus `json:"status"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Encrypted bool               `json:"encrypted,omitempty"`
	UpdateVer int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *BasicAuthResp) CopyPassword(field entity.EncryptedField) error {
	resp.Password = field.String()
	return nil
}

func TransformBasicAuth(setting *entity.Setting) (resp *BasicAuthResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config := setting.MustAsBasicAuth()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.Password.IsEncrypted()
	if resp.Encrypted {
		resp.Password = maskedPassword
	}
	return resp, nil
}

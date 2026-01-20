package sshkeydto

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
	maskedSecretKey = "****************"
)

type GetSSHKeyReq struct {
	settings.GetSettingReq
}

func NewGetSSHKeyReq() *GetSSHKeyReq {
	return &GetSSHKeyReq{}
}

func (req *GetSSHKeyReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSSHKeyResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *SSHKeyResp       `json:"data"`
}

type SSHKeyResp struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Status     base.SettingStatus `json:"status"`
	PrivateKey string             `json:"privateKey"`
	Passphrase string             `json:"passphrase,omitempty"`
	Encrypted  bool               `json:"encrypted,omitempty"`
	UpdateVer  int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func (resp *SSHKeyResp) CopyPrivateKey(field entity.EncryptedField) error {
	resp.PrivateKey = field.String()
	return nil
}

func (resp *SSHKeyResp) CopyPassphrase(field entity.EncryptedField) error {
	resp.Passphrase = field.String()
	return nil
}

func TransformSSHKey(setting *entity.Setting) (resp *SSHKeyResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	sshKey := setting.MustAsSSHKey()
	if err = copier.Copy(&resp, &sshKey); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = sshKey.PrivateKey.IsEncrypted()
	if resp.Encrypted {
		resp.PrivateKey = maskedSecretKey
		if resp.Passphrase != "" {
			resp.Passphrase = maskedSecretKey
		}
	}
	return resp, nil
}

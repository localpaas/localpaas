package ssldto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedKey = "****************"
)

type GetSslReq struct {
	settings.GetSettingReq
}

func NewGetSslReq() *GetSslReq {
	return &GetSslReq{}
}

func (req *GetSslReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSslResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *SslResp      `json:"data"`
}

type SslResp struct {
	*settings.BaseSettingResp
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
	KeySize     int    `json:"keySize"`
	Provider    string `json:"provider"`
	Email       string `json:"email"`
	Encrypted   bool   `json:"encrypted,omitempty"`
}

func (resp *SslResp) CopyPrivateKey(field entity.EncryptedField) error {
	resp.PrivateKey = field.String()
	return nil
}

func TransformSsl(setting *entity.Setting) (resp *SslResp, err error) {
	config := setting.MustAsSsl()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.PrivateKey.IsEncrypted()
	if resp.Encrypted {
		resp.PrivateKey = maskedKey
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

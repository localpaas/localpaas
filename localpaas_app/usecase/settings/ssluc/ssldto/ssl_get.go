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
	maskedSecret = "****************"
)

type GetSSLReq struct {
	settings.GetSettingReq
}

func NewGetSSLReq() *GetSSLReq {
	return &GetSSLReq{}
}

func (req *GetSSLReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSSLResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *SSLResp      `json:"data"`
}

type SSLResp struct {
	*settings.BaseSettingResp
	Certificate  string `json:"certificate"`
	PrivateKey   string `json:"privateKey"`
	KeySize      int    `json:"keySize"`
	Provider     string `json:"provider"`
	Email        string `json:"email"`
	SecretMasked bool   `json:"secretMasked,omitempty"`
}

func (resp *SSLResp) CopyPrivateKey(field entity.EncryptedField) error {
	resp.PrivateKey = field.String()
	return nil
}

func TransformSSL(
	setting *entity.Setting,
	_ *entity.RefObjects,
) (resp *SSLResp, err error) {
	config := setting.MustAsSSL()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.SecretMasked = config.PrivateKey.IsEncrypted()
	if resp.SecretMasked {
		resp.PrivateKey = maskedSecret
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

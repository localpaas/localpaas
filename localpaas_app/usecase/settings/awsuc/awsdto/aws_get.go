package awsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

const (
	maskedSecretKey = "****************"
)

type GetAWSReq struct {
	settings.GetSettingReq
}

func NewGetAWSReq() *GetAWSReq {
	return &GetAWSReq{}
}

func (req *GetAWSReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetAWSResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *AWSResp      `json:"data"`
}

type AWSResp struct {
	*settings.BaseSettingResp
	Kind        string `json:"kind,omitempty"`
	AccessKeyID string `json:"accessKeyId"`
	SecretKey   string `json:"secretKey,omitempty"`
	Region      string `json:"region"`
	Encrypted   bool   `json:"encrypted,omitempty"`
}

func (resp *AWSResp) CopySecretKey(field entity.EncryptedField) error {
	resp.SecretKey = field.String()
	return nil
}

func TransformAWS(setting *entity.Setting) (resp *AWSResp, err error) {
	config := setting.MustAsAWS()
	if err = copier.Copy(&resp, &config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Encrypted = config.SecretKey.IsEncrypted()
	if resp.Encrypted {
		resp.SecretKey = maskedSecretKey
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

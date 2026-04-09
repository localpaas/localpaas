package sslcertsettingsdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetUniqueSSLCertSettingsReq struct {
	settings.GetUniqueSettingReq
}

func NewGetUniqueSSLCertSettingsReq() *GetUniqueSSLCertSettingsReq {
	return &GetUniqueSSLCertSettingsReq{}
}

func (req *GetUniqueSSLCertSettingsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetUniqueSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetUniqueSSLCertSettingsResp struct {
	Meta *basedto.Meta        `json:"meta"`
	Data *SSLCertSettingsResp `json:"data"`
}

type SSLCertSettingsResp struct {
	*settings.BaseSettingResp
	CertType    base.SSLCertType  `json:"certType"`
	KeyType     base.SSLKeyType   `json:"keyType"`
	ValidPeriod timeutil.Duration `json:"validPeriod"`
	RootDomain  string            `json:"rootDomain"`
	Email       string            `json:"email"`
	AutoRenew   bool              `json:"autoRenew"`
}

func TransformSSLCertSettings(
	setting *entity.Setting,
	_ *entity.RefObjects,
) (resp *SSLCertSettingsResp, err error) {
	config := setting.MustAsSSLCertSettings()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

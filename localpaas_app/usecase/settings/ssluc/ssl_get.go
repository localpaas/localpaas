package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

func (uc *SSLUC) GetSSL(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.GetSSLReq,
) (*ssldto.GetSSLResp, error) {
	req.Type = currentSettingType
	setting, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsSSL().MustDecrypt()
	resp, err := ssldto.TransformSSL(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.GetSSLResp{
		Data: resp,
	}, nil
}

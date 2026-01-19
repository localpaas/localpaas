package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

func (uc *SslUC) GetSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.GetSslReq,
) (*ssldto.GetSslResp, error) {
	req.Type = currentSettingType
	setting, err := providers.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &providers.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsSsl().MustDecrypt()
	resp, err := ssldto.TransformSsl(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.GetSslResp{
		Data: resp,
	}, nil
}

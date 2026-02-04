package accesstokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
)

func (uc *AccessTokenUC) GetAccessToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *accesstokendto.GetAccessTokenReq,
) (*accesstokendto.GetAccessTokenResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsAccessToken().MustDecrypt()
	resp, err := accesstokendto.TransformAccessToken(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &accesstokendto.GetAccessTokenResp{
		Data: resp,
	}, nil
}

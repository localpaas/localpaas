package accesstokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
)

func (uc *AccessTokenUC) ListAccessToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *accesstokendto.ListAccessTokenReq,
) (*accesstokendto.ListAccessTokenResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := accesstokendto.TransformAccessTokens(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &accesstokendto.ListAccessTokenResp{
		Data: respData,
	}, nil
}

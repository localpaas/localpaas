package accesstokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
)

func (uc *AccessTokenUC) UpdateAccessTokenMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *accesstokendto.UpdateAccessTokenMetaReq,
) (*accesstokendto.UpdateAccessTokenMetaResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &accesstokendto.UpdateAccessTokenMetaResp{}, nil
}

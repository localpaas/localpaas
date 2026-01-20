package basicauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) UpdateBasicAuthMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.UpdateBasicAuthMetaReq,
) (*basicauthdto.UpdateBasicAuthMetaResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &basicauthdto.UpdateBasicAuthMetaResp{}, nil
}

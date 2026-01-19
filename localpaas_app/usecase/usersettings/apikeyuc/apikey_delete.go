package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) DeleteAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.DeleteAPIKeyReq,
) (*apikeydto.DeleteAPIKeyResp, error) {
	req.Type = currentSettingType
	_, err := providers.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &providers.DeleteSettingData{
		SettingRepo: uc.settingRepo,
		ExtraLoadOpts: []bunex.SelectQueryOption{
			bunex.SelectWhere("setting.object_id = ?", auth.User.ID),
		},
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.DeleteAPIKeyResp{}, nil
}

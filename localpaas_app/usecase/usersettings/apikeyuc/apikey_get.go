package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) GetAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.GetAPIKeyReq,
) (*apikeydto.GetAPIKeyResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := apikeydto.TransformAPIKey(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.GetAPIKeyResp{
		Data: resp,
	}, nil
}

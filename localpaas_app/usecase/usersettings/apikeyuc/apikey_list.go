package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) ListAPIKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.ListAPIKeyReq,
) (*apikeydto.ListAPIKeyResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := apikeydto.TransformAPIKeys(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.ListAPIKeyResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

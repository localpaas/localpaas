package apikeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

func (uc *APIKeyUC) UpdateAPIKeyMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *apikeydto.UpdateAPIKeyMetaReq,
) (*apikeydto.UpdateAPIKeyMetaResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &providers.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &apikeydto.UpdateAPIKeyMetaResp{}, nil
}

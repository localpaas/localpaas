package basicauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) GetBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.GetBasicAuthReq,
) (*basicauthdto.GetBasicAuthResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsBasicAuth().MustDecrypt()
	resp, err := basicauthdto.TransformBasicAuth(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &basicauthdto.GetBasicAuthResp{
		Data: resp,
	}, nil
}

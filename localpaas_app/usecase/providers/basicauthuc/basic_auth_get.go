package basicauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) GetBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.GetBasicAuthReq,
) (*basicauthdto.GetBasicAuthResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeBasicAuth, req.ID, false)
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

package basicauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/basicauthuc/basicauthdto"
)

func (uc *BasicAuthUC) ListBasicAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *basicauthdto.ListBasicAuthReq,
) (*basicauthdto.ListBasicAuthResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := basicauthdto.TransformBasicAuths(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &basicauthdto.ListBasicAuthResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

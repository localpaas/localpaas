package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) ListRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.ListRegistryAuthReq,
) (*registryauthdto.ListRegistryAuthResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := registryauthdto.TransformRegistryAuths(resp.Data, req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.ListRegistryAuthResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

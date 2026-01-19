package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) ListRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.ListRegistryAuthReq,
) (*registryauthdto.ListRegistryAuthResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := registryauthdto.TransformRegistryAuths(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.ListRegistryAuthResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

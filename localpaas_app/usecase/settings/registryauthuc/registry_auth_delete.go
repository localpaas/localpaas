package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) DeleteRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.DeleteRegistryAuthReq,
) (*registryauthdto.DeleteRegistryAuthResp, error) {
	req.Type = currentSettingType
	_, err := settings.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &settings.DeleteSettingData{
		SettingRepo:              uc.settingRepo,
		ProjectSharedSettingRepo: uc.projectSharedSettingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.DeleteRegistryAuthResp{}, nil
}

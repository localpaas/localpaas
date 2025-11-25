package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) GetRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.GetRegistryAuthReq,
) (*registryauthdto.GetRegistryAuthResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeRegistryAuth, req.ID, false)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsRegistryAuth().MustDecrypt()
	resp, err := registryauthdto.TransformRegistryAuth(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.GetRegistryAuthResp{
		Data: resp,
	}, nil
}

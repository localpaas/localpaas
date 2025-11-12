package registryauthuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

func (uc *RegistryAuthUC) GetRegistryAuth(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.GetRegistryAuthReq,
) (*registryauthdto.GetRegistryAuthResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, req.ID,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeRegistryAuth),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := registryauthdto.TransformRegistryAuth(setting, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &registryauthdto.GetRegistryAuthResp{
		Data: resp,
	}, nil
}

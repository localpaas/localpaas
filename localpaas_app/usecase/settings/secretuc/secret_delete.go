package secretuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

func (uc *SecretUC) DeleteSecret(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.DeleteSecretReq,
) (*secretdto.DeleteSecretResp, error) {
	req.Type = currentSettingType
	_, err := settings.DeleteSetting(ctx, uc.db, &req.DeleteSettingReq, &settings.DeleteSettingData{
		SettingRepo:              uc.settingRepo,
		ProjectSharedSettingRepo: uc.projectSharedSettingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &secretdto.DeleteSecretResp{}, nil
}

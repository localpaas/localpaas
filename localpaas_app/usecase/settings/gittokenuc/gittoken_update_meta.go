package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) UpdateGitTokenMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.UpdateGitTokenMetaReq,
) (*gittokendto.UpdateGitTokenMetaResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.UpdateGitTokenMetaResp{}, nil
}

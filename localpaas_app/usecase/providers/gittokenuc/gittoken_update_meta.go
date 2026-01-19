package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) UpdateGitTokenMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.UpdateGitTokenMetaReq,
) (*gittokendto.UpdateGitTokenMetaResp, error) {
	req.Type = currentSettingType
	_, err := providers.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &providers.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.UpdateGitTokenMetaResp{}, nil
}

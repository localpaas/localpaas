package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) GetGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.GetGitTokenReq,
) (*gittokendto.GetGitTokenResp, error) {
	req.Type = currentSettingType
	setting, err := providers.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &providers.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsGitToken().MustDecrypt()
	resp, err := gittokendto.TransformGitToken(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.GetGitTokenResp{
		Data: resp,
	}, nil
}

package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) GetGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.GetGitTokenReq,
) (*gittokendto.GetGitTokenResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsGitToken().MustDecrypt()
	resp, err := gittokendto.TransformGitToken(setting, req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.GetGitTokenResp{
		Data: resp,
	}, nil
}

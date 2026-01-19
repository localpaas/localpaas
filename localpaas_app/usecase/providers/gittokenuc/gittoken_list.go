package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) ListGitToken(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.ListGitTokenReq,
) (*gittokendto.ListGitTokenResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := gittokendto.TransformGitTokens(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.ListGitTokenResp{
		Data: respData,
	}, nil
}

package githubappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
)

func (uc *GithubAppUC) GetGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.GetGithubAppReq,
) (*githubappdto.GetGithubAppResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, base.SettingTypeGithubApp, req.ID, false)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsGithubApp().MustDecrypt()
	resp, err := githubappdto.TransformGithubApp(setting, config.Current.SsoBaseCallbackURL())
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.GetGithubAppResp{
		Data: resp,
	}, nil
}

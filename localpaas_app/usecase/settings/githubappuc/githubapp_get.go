package githubappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
)

func (uc *GithubAppUC) GetGithubApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.GetGithubAppReq,
) (*githubappdto.GetGithubAppResp, error) {
	req.Type = currentSettingType
	setting, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsGithubApp().MustDecrypt()
	resp, err := githubappdto.TransformGithubApp(setting, config.Current.SsoBaseCallbackURL(), req.ObjectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.GetGithubAppResp{
		Data: resp,
	}, nil
}

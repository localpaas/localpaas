package githubappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
	"github.com/localpaas/localpaas/services/github"
)

func (uc *GithubAppUC) TestGithubAppConn(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.TestGithubAppConnReq,
) (*githubappdto.TestGithubAppConnResp, error) {
	app, err := github.NewApp(req.AppID, req.InstallationID, reflectutil.UnsafeStrToBytes(req.PrivateKey))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	_, _, err = app.ListInstallations(ctx, &basedto.Paging{Limit: 1})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.TestGithubAppConnResp{}, nil
}

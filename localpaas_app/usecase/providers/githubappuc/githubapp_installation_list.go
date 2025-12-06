package githubappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/githubappuc/githubappdto"
	"github.com/localpaas/localpaas/services/github"
)

func (uc *GithubAppUC) ListAppInstallation(
	ctx context.Context,
	auth *basedto.Auth,
	req *githubappdto.ListAppInstallationReq,
) (*githubappdto.ListAppInstallationResp, error) {
	app, err := github.NewApp(req.AppID, req.InstallationID, reflectutil.UnsafeStrToBytes(req.PrivateKey))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	installations, pagingMeta, err := app.ListInstallations(ctx, &req.Paging)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := githubappdto.TransformAppInstallations(installations)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &githubappdto.ListAppInstallationResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: resp,
	}, nil
}

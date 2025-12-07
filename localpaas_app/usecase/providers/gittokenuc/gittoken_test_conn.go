package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
	"github.com/localpaas/localpaas/services/github"
)

func (uc *GitTokenUC) TestGitTokenConn(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.TestGitTokenConnReq,
) (*gittokendto.TestGitTokenConnResp, error) {
	var err error
	switch req.Kind { //nolint:exhaustive
	case base.GitSourceGithub:
		err = uc.testGithubTokenConn(ctx, req)
	case base.GitSourceGitlab:
		err = uc.testGitlabTokenConn(ctx, req)
	case base.GitSourceGitea:
		err = uc.testGiteaTokenConn(ctx, req)
	default:
		err = apperrors.New(apperrors.ErrUnsupported).
			WithMsgLog("Git source '%s' unsupported", req.Kind)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gittokendto.TestGitTokenConnResp{}, nil
}

func (uc *GitTokenUC) testGithubTokenConn(
	ctx context.Context,
	req *gittokendto.TestGitTokenConnReq,
) error {
	client, err := github.NewFromPersonalToken(req.Token)
	if err != nil {
		return apperrors.Wrap(err)
	}
	_, err = client.ListRepos(ctx, github.ListOptionPerPage(1))
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *GitTokenUC) testGitlabTokenConn(
	_ context.Context,
	_ *gittokendto.TestGitTokenConnReq,
) error {
	// TODO: add implementation
	return apperrors.Wrap(apperrors.ErrNotImplemented)
}

func (uc *GitTokenUC) testGiteaTokenConn(
	_ context.Context,
	_ *gittokendto.TestGitTokenConnReq,
) error {
	// TODO: add implementation
	return apperrors.Wrap(apperrors.ErrNotImplemented)
}

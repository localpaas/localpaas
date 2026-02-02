package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/gittokenuc/gittokendto"
	"github.com/localpaas/localpaas/services/gitea"
	"github.com/localpaas/localpaas/services/github"
	"github.com/localpaas/localpaas/services/gitlab"
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
	case base.GitSourceGitlab, base.GitSourceGitlabCustom:
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
	_, _, err = client.ListUserRepos(ctx, &basedto.Paging{Limit: 1})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *GitTokenUC) testGitlabTokenConn(
	ctx context.Context,
	req *gittokendto.TestGitTokenConnReq,
) error {
	client, err := gitlab.NewFromToken(req.Token, req.BaseURL)
	if err != nil {
		return apperrors.Wrap(err)
	}
	_, _, err = client.ListAllProjects(ctx, &basedto.Paging{Limit: 1})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *GitTokenUC) testGiteaTokenConn(
	ctx context.Context,
	req *gittokendto.TestGitTokenConnReq,
) error {
	client, err := gitea.NewFromToken(req.Token, req.BaseURL)
	if err != nil {
		return apperrors.Wrap(err)
	}
	_, _, err = client.ListAllRepos(ctx, &basedto.Paging{Limit: 1})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

package gitsourceuc

import (
	"context"

	gogithub "github.com/google/go-github/v75/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/gitsourceuc/gitsourcedto"
	"github.com/localpaas/localpaas/services/github"
)

func (uc *GitSourceUC) ListRepo(
	ctx context.Context,
	auth *basedto.Auth,
	req *gitsourcedto.ListRepoReq,
) (*gitsourcedto.ListRepoResp, error) {
	setting, err := uc.settingRepo.GetByID(ctx, uc.db, "", req.SettingID, true,
		bunex.SelectWhere("setting.type IN (?)",
			bunex.InItems(base.SettingTypeGithubApp, base.SettingTypeGitToken)),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	switch setting.Type { //nolint:exhaustive
	case base.SettingTypeGithubApp:
		return uc.listGithubRepo(ctx, req, setting)
	case base.SettingTypeGitToken:
		switch base.GitSource(setting.Kind) { //nolint:exhaustive
		case base.GitSourceGithub:
			return uc.listGithubRepo(ctx, req, setting)
		case base.GitSourceGitlab:
			return uc.listGitlabRepo(ctx, req, setting)
		case base.GitSourceGitea:
			return uc.listGiteaRepo(ctx, req, setting)
		}
	}

	return nil, apperrors.NewUnsupported()
}

func (uc *GitSourceUC) listGithubRepo(
	ctx context.Context,
	req *gitsourcedto.ListRepoReq,
	setting *entity.Setting,
) (*gitsourcedto.ListRepoResp, error) {
	client, err := github.NewFromSetting(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	var repos []*gogithub.Repository
	var pagingMeta *basedto.PagingMeta
	if req.Paging.Limit > github.MaxListPageSize {
		repos, err = client.ListAllRepos(ctx)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		pagingMeta = &basedto.PagingMeta{
			Total: len(repos),
		}
	} else {
		repos, pagingMeta, err = client.ListRepos(ctx, &req.Paging)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	resp, err := gitsourcedto.TransformGithubRepos(repos)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gitsourcedto.ListRepoResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: resp,
	}, nil
}

func (uc *GitSourceUC) listGitlabRepo(
	_ context.Context,
	_ *gitsourcedto.ListRepoReq,
	_ *entity.Setting,
) (*gitsourcedto.ListRepoResp, error) {
	// TODO: add implementation
	return nil, apperrors.Wrap(apperrors.ErrNotImplemented)
}

func (uc *GitSourceUC) listGiteaRepo(
	_ context.Context,
	_ *gitsourcedto.ListRepoReq,
	_ *entity.Setting,
) (*gitsourcedto.ListRepoResp, error) {
	// TODO: add implementation
	return nil, apperrors.Wrap(apperrors.ErrNotImplemented)
}

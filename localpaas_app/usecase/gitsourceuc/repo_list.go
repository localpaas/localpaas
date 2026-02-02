package gitsourceuc

import (
	"context"

	gogithub "github.com/google/go-github/v79/github"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/gitsourceuc/gitsourcedto"
	"github.com/localpaas/localpaas/services/gitea"
	"github.com/localpaas/localpaas/services/github"
	"github.com/localpaas/localpaas/services/gitlab"
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
		case base.GitSourceGitlab, base.GitSourceGitlabCustom:
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
	if client.IsAppClient() {
		repos, pagingMeta, err = client.ListAppRepos(ctx, &req.Paging)
	} else {
		repos, pagingMeta, err = client.ListUserRepos(ctx, &req.Paging)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := gitsourcedto.TransformGithubRepos(repos)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gitsourcedto.ListRepoResp{
		Meta: &basedto.ListMeta{Page: pagingMeta},
		Data: resp,
	}, nil
}

func (uc *GitSourceUC) listGitlabRepo(
	ctx context.Context,
	req *gitsourcedto.ListRepoReq,
	setting *entity.Setting,
) (*gitsourcedto.ListRepoResp, error) {
	client, err := gitlab.NewFromSetting(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	projects, pagingMeta, err := client.ListProjects(ctx, &req.Paging)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := gitsourcedto.TransformGitlabProjects(projects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gitsourcedto.ListRepoResp{
		Meta: &basedto.ListMeta{Page: pagingMeta},
		Data: resp,
	}, nil
}

func (uc *GitSourceUC) listGiteaRepo(
	ctx context.Context,
	req *gitsourcedto.ListRepoReq,
	setting *entity.Setting,
) (*gitsourcedto.ListRepoResp, error) {
	client, err := gitea.NewFromSetting(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	repos, pagingMeta, err := client.ListRepos(ctx, &req.Paging)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := gitsourcedto.TransformGiteaRepos(repos)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &gitsourcedto.ListRepoResp{
		Meta: &basedto.ListMeta{Page: pagingMeta},
		Data: resp,
	}, nil
}

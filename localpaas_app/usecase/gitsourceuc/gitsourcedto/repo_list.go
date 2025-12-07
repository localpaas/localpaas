package gitsourcedto

import (
	"strconv"

	"github.com/google/go-github/v75/github"
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type ListRepoReq struct {
	SettingID string         `json:"-"`
	Paging    basedto.Paging `json:"-"`
}

func NewListRepoReq() *ListRepoReq {
	return &ListRepoReq{}
}

func (req *ListRepoReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type ListRepoResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*RepoResp   `json:"data"`
}

type RepoResp struct {
	ID            string `json:"id" copy:"-"`
	Name          string `json:"name"`
	FullName      string `json:"fullName"`
	DefaultBranch string `json:"defaultBranch"`
	MasterBranch  string `json:"masterBranch"`
	CloneURL      string `json:"cloneURL"`
	GitURL        string `json:"gitURL"`
}

func TransformGithubRepo(repo *github.Repository) (resp *RepoResp, err error) {
	if err = copier.Copy(&resp, repo); err != nil {
		return nil, apperrors.Wrap(err)
	}
	resp.ID = strconv.FormatInt(*repo.ID, 10)
	return resp, nil
}

func TransformGithubRepos(repos []*github.Repository) ([]*RepoResp, error) {
	resp, err := basedto.TransformObjectSlice(repos, TransformGithubRepo)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

package repocheckoutservice

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
)

type RepoCheckoutReq struct {
	Project     *entity.Project
	App         *entity.App
	RepoSource  *entity.DeploymentRepoSource
	NoCache     bool
	CredSetting *entity.Setting

	RefObjects  *entity.RefObjects
	LogStore    *tasklog.Store
	TempDir     string
	CheckoutDir string
}

type RepoCheckoutResp struct {
	CommitHash    string
	CommitMessage string
	CommitTitle   string
	CommitAuthor  string
}

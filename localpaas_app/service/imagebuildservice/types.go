package imagebuildservice

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
)

type ImageBuildReq struct {
	Project            *entity.Project
	App                *entity.App
	RepoSource         *entity.DeploymentRepoSource
	ImageBuildSettings *entity.ImageBuildSettings
	NoCache            bool

	BuildID     string
	RefObjects  *entity.RefObjects
	LogStore    *tasklog.Store
	TempDir     string
	CheckoutDir string
}

type ImageBuildResp struct {
	ImageTags []string
}

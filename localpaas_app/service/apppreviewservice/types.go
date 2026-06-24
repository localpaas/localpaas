package apppreviewservice

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type CreatePreviewReq struct {
	ProjectID   string
	AppID       string
	PullRequest string

	OnInitDeployment func(*entity.Deployment) error
	OnDeploymentTask func(*entity.Task) error
}

type CreatePreviewResp struct {
	PreviewApp     *entity.App
	Deployment     *entity.Deployment
	DeploymentTask *entity.Task
	OnCleanup      func(error) error
}

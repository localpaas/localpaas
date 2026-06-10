package appdeploymentservice

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/tasks/queue"
)

type AppDeploymentReq struct {
	*queue.TaskExecData
}

type AppDeploymentResp struct {
	Deployment *entity.Deployment
}

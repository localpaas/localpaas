package appdeploymentservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	CreateDeploymentAndTask(app *entity.App, deploymentSettings *entity.AppDeploymentSettings) (
		*entity.Deployment, *entity.Task, error)

	Deploy(ctx context.Context, db database.Tx, req *AppDeploymentReq) (*AppDeploymentResp, error)
}

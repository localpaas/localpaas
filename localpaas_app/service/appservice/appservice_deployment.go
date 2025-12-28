package appservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

func (s *appService) CancelAllDeployments(ctx context.Context, db database.Tx, app *entity.App) error {
	// Cancel all not-started deployments in the DB
	deployments, _, err := s.deploymentRepo.List(ctx, db, app.ID, nil,
		bunex.SelectWhere("deployment.status = ?", base.DeploymentStatusNotStarted),
		bunex.SelectFor("UPDATE OF deployment SKIP LOCKED"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	timeNow := timeutil.NowUTC()
	for _, deployment := range deployments {
		deployment.Status = base.DeploymentStatusCanceled
		deployment.UpdatedAt = timeNow
	}
	err = s.deploymentRepo.UpsertMulti(ctx, db, deployments,
		entity.DeploymentUpsertingConflictCols, []string{"status", "updated_at"})
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Cancel all in-progress deployments
	err = s.deploymentInfoRepo.CancelAllOfApp(ctx, app.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

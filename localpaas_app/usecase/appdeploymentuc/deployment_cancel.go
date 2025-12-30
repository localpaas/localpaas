package appdeploymentuc

import (
	"context"
	"errors"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc/appdeploymentdto"
)

func (uc *AppDeploymentUC) CancelDeployment(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdeploymentdto.CancelDeploymentReq,
) (*appdeploymentdto.CancelDeploymentResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		deployment, err := uc.deploymentRepo.GetByID(ctx, db, req.AppID, req.DeploymentID,
			bunex.SelectFor("UPDATE OF deployment SKIP LOCKED"),
		)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.Wrap(err)
		}

		if deployment != nil {
			if deployment.Status == base.DeploymentStatusDone ||
				deployment.Status == base.DeploymentStatusFailed ||
				deployment.Status == base.DeploymentStatusCanceled {
				return apperrors.New(apperrors.ErrStatusNotAllowAction)
			}

			deployment.Status = base.DeploymentStatusCanceled
			deployment.UpdatedAt = timeutil.NowUTC()

			err = uc.deploymentRepo.Update(ctx, db, deployment,
				bunex.UpdateColumns("status", "updated_at"),
			)
			if err != nil {
				return apperrors.Wrap(err)
			}
			return nil
		}

		// Deployment is in-progress, set `cancel` flag of the deployment info in redis
		deploymentInfo, err := uc.deploymentInfoRepo.Get(ctx, req.DeploymentID)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				return apperrors.New(apperrors.ErrInternalServer).
					WithMsgLog("deployment info not found, please try again later")
			}
			return apperrors.Wrap(err)
		}

		deploymentInfo.Cancel = true
		err = uc.deploymentInfoRepo.Update(ctx, req.DeploymentID, deploymentInfo)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdeploymentdto.CancelDeploymentResp{}, nil
}

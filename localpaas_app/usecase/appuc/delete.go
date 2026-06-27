package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *UC) DeleteApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.DeleteAppReq,
) (*appdto.DeleteAppResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
			bunex.SelectFor("UPDATE OF app"),
		)
		if err != nil {
			return apperrors.New(err)
		}

		// Remove app and its data from the infra
		err = uc.appService.DeleteApp(ctx, db, app)
		if err != nil {
			return apperrors.New(err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &appdto.DeleteAppResp{}, nil
}

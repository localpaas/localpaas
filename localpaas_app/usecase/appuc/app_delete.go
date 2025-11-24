package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *AppUC) DeleteApp(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.DeleteAppReq,
) (*appdto.DeleteAppResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData := &deleteAppData{}
		err := uc.loadAppDataForDelete(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingAppData{}
		uc.prepareDeletingApp(appData, persistingData)

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		// Remove service for the app in docker
		err = uc.dockerManager.ServiceRemove(ctx, appData.App.ServiceID)
		if err != nil {
			return apperrors.Wrap(err)
		}

		// Remove app config from nginx
		err = uc.nginxService.RemoveAppConfig(ctx, appData.App)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.DeleteAppResp{}, nil
}

type deleteAppData struct {
	App *entity.App
}

func (uc *AppUC) loadAppDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *appdto.DeleteAppReq,
	data *deleteAppData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app

	if app.Status == base.AppStatusDeleting { //nolint
		// TODO: handle task deletion if previously failed
	}

	return nil
}

func (uc *AppUC) prepareDeletingApp(
	data *deleteAppData,
	persistingData *persistingAppData,
) {
	app := data.App
	app.Status = base.AppStatusDeleting
	persistingData.UpsertingApps = append(persistingData.UpsertingApps, app)
}

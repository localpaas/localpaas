package appuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

func (uc *UC) UpdateAppStatus(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppStatusReq,
) (*appdto.UpdateAppStatusResp, error) {
	var oldAppStatus base.AppStatus
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData := &updateAppData{}
		err := uc.loadAppDataForUpdateStatus(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}
		if !appData.HasChanges {
			return nil
		}

		oldAppStatus = appData.App.Status
		persistingData := &persistingAppData{}
		uc.preparePersistingAppStatusUpdate(req, appData, persistingData)

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.appService.OnAppStatusChanged(ctx, appData.App, oldAppStatus)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppStatusResp{}, nil
}

func (uc *UC) loadAppDataForUpdateStatus(
	ctx context.Context,
	db database.IDB,
	req *appdto.UpdateAppStatusReq,
	data *updateAppData,
) error {
	app, err := uc.appService.LoadApp(ctx, db, req.ProjectID, req.ID, true, false,
		bunex.SelectFor("UPDATE OF app"),
		bunex.SelectRelation("Project"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if app.UpdateVer != req.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.App = app
	data.HasChanges = app.Status != req.Status

	return nil
}

func (uc *UC) preparePersistingAppStatusUpdate(
	req *appdto.UpdateAppStatusReq,
	data *updateAppData,
	persistingData *persistingAppData,
) {
	timeNow := timeutil.NowUTC()
	app := data.App
	app.UpdateVer++
	app.Status = req.Status
	app.UpdatedAt = timeNow

	persistingData.UpsertingApps = append(persistingData.UpsertingApps, app)
}

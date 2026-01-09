package appuc

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
	"github.com/localpaas/localpaas/services/docker"
)

func (uc *AppUC) UpdateAppServiceSpec(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.UpdateAppServiceSpecReq,
) (*appdto.UpdateAppServiceSpecResp, error) {
	var data *updateAppServiceSpecData
	var persistingData *persistingAppData
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		data = &updateAppServiceSpecData{}
		err := uc.loadAppServiceSpecForUpdate(ctx, db, req, data)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData = &persistingAppData{}
		uc.prepareUpdatingAppServiceSpec(req, data)

		err = uc.persistData(ctx, db, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.applyAppServiceSpec(ctx, data)
		if err != nil {
			return apperrors.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.UpdateAppServiceSpecResp{}, nil
}

type updateAppServiceSpecData struct {
	App      *entity.App
	Service  *swarm.Service
	Errors   []string // stores errors
	Warnings []string // stores warnings
}

func (uc *AppUC) loadAppServiceSpecForUpdate(
	ctx context.Context,
	db database.Tx,
	req *appdto.UpdateAppServiceSpecReq,
	data *updateAppServiceSpecData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE OF app"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.App = app

	service, err := uc.dockerManager.ServiceInspect(ctx, app.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Service = service

	if data.Service == nil || data.Service.Version.Index != uint64(req.UpdateVer) { //nolint:gosec
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}

	return nil
}

func (uc *AppUC) prepareUpdatingAppServiceSpec(
	req *appdto.UpdateAppServiceSpecReq,
	data *updateAppServiceSpecData,
) {
	service := data.Service
	spec := &service.Spec

	docker.ApplyServiceModeSpec(spec, req.ServiceMode)
	docker.ApplyServiceEndpointSpec(spec, req.EndpointSpec)
	docker.ApplyServiceTaskSpec(&spec.TaskTemplate, req.TaskSpec)
	docker.ApplyServiceContainerSpec(spec.TaskTemplate.ContainerSpec, req.ContainerSpec)
}

func (uc *AppUC) applyAppServiceSpec(
	ctx context.Context,
	data *updateAppServiceSpecData,
) error {
	service := data.Service

	_, err := uc.dockerManager.ServiceUpdate(ctx, service.ID, &service.Version, &service.Spec)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

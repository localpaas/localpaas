package appsettingsuc

import (
	"context"

	"github.com/docker/docker/api/types/mount"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appsettingsuc/appsettingsdto"
)

func (uc *UC) GetAppStorageSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *appsettingsdto.GetAppStorageSettingsReq,
) (*appsettingsdto.GetAppStorageSettingsResp, error) {
	app, err := uc.appRepo.GetByID(ctx, uc.db, req.ProjectID, req.AppID,
		bunex.SelectExcludeColumns(entity.AppDefaultExcludeColumns...),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	service, err := uc.appService.ServiceInspect(ctx, app.ServiceID, true)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Filter out unsupported mount types
	returningMounts := make([]mount.Mount, 0, len(service.Spec.TaskTemplate.ContainerSpec.Mounts))
	for _, mnt := range service.Spec.TaskTemplate.ContainerSpec.Mounts {
		if gofn.Contain(supportedMountTypes, mnt.Type) {
			returningMounts = append(returningMounts, mnt)
		}
	}
	service.Spec.TaskTemplate.ContainerSpec.Mounts = returningMounts

	resp, err := appsettingsdto.TransformStorageSettings(service)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appsettingsdto.GetAppStorageSettingsResp{
		Data: resp,
	}, nil
}

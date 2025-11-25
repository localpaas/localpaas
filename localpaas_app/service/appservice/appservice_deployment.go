package appservice

import (
	"context"

	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type AppDeploymentReq struct {
	Deployment              *entity.Setting
	ImageSourceRegistryAuth *entity.Setting
}

type AppDeploymentResp struct {
}

func (s *appService) UpdateAppDeployment(ctx context.Context, app *entity.App, req *AppDeploymentReq) (
	*AppDeploymentResp, error) {
	deploymentSettings := req.Deployment.MustAsAppDeploymentSettings()
	switch {
	case deploymentSettings.ImageSource != nil && deploymentSettings.ImageSource.Enabled:
		return s.updateAppDeploymentImageSource(ctx, app, req)
	case deploymentSettings.CodeSource != nil && deploymentSettings.CodeSource.Enabled:
		return s.updateAppDeploymentImageSource(ctx, app, req)
	}
	return nil, nil
}

func (s *appService) updateAppDeploymentImageSource(ctx context.Context, app *entity.App, req *AppDeploymentReq) (
	*AppDeploymentResp, error) {
	imageSource := req.Deployment.MustAsAppDeploymentSettings().ImageSource

	service, err := s.dockerManager.ServiceInspect(ctx, app.ServiceID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	spec := &service.Spec
	spec.TaskTemplate.ContainerSpec.Image = imageSource.Name

	var regAuthHeader string
	if req.ImageSourceRegistryAuth != nil {
		regAuthHeader, err = req.ImageSourceRegistryAuth.MustAsRegistryAuth().MustDecrypt().GenerateAuthHeader()
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	_, err = s.dockerManager.ServiceUpdate(ctx, app.ServiceID, &service.Version, spec,
		func(options *swarm.ServiceUpdateOptions) {
			options.EncodedRegistryAuth = regAuthHeader
		})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &AppDeploymentResp{}, nil
}

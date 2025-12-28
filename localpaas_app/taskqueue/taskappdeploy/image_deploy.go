package taskappdeploy

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/services/docker"
)

type imageDeployTaskData struct {
	*taskData
	RegistryAuth *entity.RegistryAuth
}

func (e *Executor) deployFromImage(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	data := &imageDeployTaskData{taskData: taskData}
	deployment := data.Deployment
	imageSource := deployment.DeploymentSettings.ImageSource

	regAuthHeader, err := e.calcRegistryAuthHeader(ctx, db, imageSource.RegistryAuth.ID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	logsReader, err := e.dockerManager.ImagePull(ctx, imageSource.Image, func(options *image.PullOptions) {
		options.RegistryAuth = regAuthHeader
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	logsChan, _ := docker.StartJSONMsgScanning(ctx, logsReader)
	for msg := range logsChan {
		// print(">>>>>>>>>> ", reflectutil.UnsafeBytesToStr(gofn.Must(json.Marshal(msg))))
		if msg.Error != nil {
			err = errors.Join(err, msg.Error)
		}
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Deployment is canceled by user while we are processing it
	data.DeploymentCanceled, err = e.isDeploymentCanceled(ctx, deployment)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if data.DeploymentCanceled {
		return nil
	}

	service, err := e.dockerManager.ServiceInspect(ctx, deployment.App.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	spec := &service.Spec
	spec.TaskTemplate.ContainerSpec.Image = imageSource.Image

	_, err = e.dockerManager.ServiceUpdate(ctx, deployment.App.ServiceID, &service.Version, spec,
		func(options *swarm.ServiceUpdateOptions) {
			options.EncodedRegistryAuth = regAuthHeader
		})
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) calcRegistryAuthHeader(
	ctx context.Context,
	db database.Tx,
	regAuthID string,
) (string, error) {
	if regAuthID == "" {
		return "", nil
	}
	setting, err := e.settingRepo.GetByID(ctx, db, base.SettingTypeRegistryAuth, regAuthID, true)
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	regAuth, err := setting.AsRegistryAuth()
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	regAuthHeader, err := regAuth.GenerateAuthHeader()
	if err != nil {
		return "", apperrors.Wrap(err)
	}
	return regAuthHeader, nil
}

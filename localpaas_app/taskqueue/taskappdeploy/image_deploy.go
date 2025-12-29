package taskappdeploy

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/services/docker"
)

type imageDeployTaskData struct {
	*taskData
	RegistryAuthHeader string
}

func (e *Executor) deployFromImage(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	data := &imageDeployTaskData{taskData: taskData}
	err := e.deployStepPullImage(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Check if deployment is canceled by user while we are processing it
	isCanceled, err := e.checkDeploymentCanceled(ctx, data.taskData)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if isCanceled {
		return nil
	}

	err = e.deployStepUpdateService(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) deployStepPullImage(
	ctx context.Context,
	db database.Tx,
	data *imageDeployTaskData,
) error {
	imageSource := data.Deployment.Settings.ImageSource

	regAuthHeader, err := e.calcRegistryAuthHeader(ctx, db, data)
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

	return nil
}

func (e *Executor) deployStepUpdateService(
	ctx context.Context,
	db database.Tx,
	data *imageDeployTaskData,
) error {
	deployment := data.Deployment
	imageSource := deployment.Settings.ImageSource

	regAuthHeader, err := e.calcRegistryAuthHeader(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
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
	data *imageDeployTaskData,
) (string, error) {
	if data.RegistryAuthHeader != "" {
		return data.RegistryAuthHeader, nil
	}
	regAuthID := data.Deployment.Settings.ImageSource.RegistryAuth.ID
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
	data.RegistryAuthHeader = regAuthHeader
	return regAuthHeader, nil
}

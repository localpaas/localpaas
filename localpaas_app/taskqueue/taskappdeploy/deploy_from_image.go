package taskappdeploy

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/swarm"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	stepImagePull = "image-pull"
)

type imageDeployTaskData struct {
	*taskData
	RegistryAuthHeader string
	Step               string
}

func (e *Executor) deployFromImage(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	data := &imageDeployTaskData{taskData: taskData}

	// 1. Pull image from the registry
	err := e.imageDeployStepImagePull(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if data.IsCanceled() {
		return nil
	}

	// 2. Pre-deployment command execution
	err = e.deployStepExecCmd(ctx, data.taskData, true)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 3. Apply image to service
	err = e.imageDeployStepServiceApply(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// 4. Post-deployment command execution
	err = e.deployStepExecCmd(ctx, data.taskData, false)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) imageDeployStepImagePull(
	ctx context.Context,
	db database.Tx,
	data *imageDeployTaskData,
) (err error) {
	data.Step = stepImagePull
	imageSource := data.Deployment.Settings.ImageSource

	e.addStepStartLog(ctx, data.taskData, "Start pulling image...")
	defer e.addStepEndLog(ctx, data.taskData, timeutil.NowUTC(), err)

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

	logsChan, _ := docker.StartScanningJSONMsg(ctx, logsReader, batchrecvchan.Options{})
	for msgs := range logsChan {
		for _, msg := range msgs {
			// print(" >>>>>>>>>>> ", msg.String())
			if msg.Error != nil {
				err = errors.Join(err, msg.Error)
			}
		}
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) imageDeployStepServiceApply(
	ctx context.Context,
	db database.Tx,
	data *imageDeployTaskData,
) (err error) {
	data.Step = stepServiceApply
	deployment := data.Deployment
	imageSource := deployment.Settings.ImageSource

	e.addStepStartLog(ctx, data.taskData, "Applying changes to service...")
	defer e.addStepEndLog(ctx, data.taskData, timeutil.NowUTC(), err)

	regAuthHeader, err := e.calcRegistryAuthHeader(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	service, err := e.dockerManager.ServiceInspect(ctx, deployment.App.ServiceID)
	if err != nil {
		return apperrors.Wrap(err)
	}

	spec := &service.Spec
	contSpec := spec.TaskTemplate.ContainerSpec
	contSpec.Image = imageSource.Image
	if deployment.Settings.WorkingDir != nil {
		contSpec.Dir = *deployment.Settings.WorkingDir
	}
	if deployment.Settings.Command != nil {
		docker.ApplyContainerCommand(contSpec, *deployment.Settings.Command)
	}

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

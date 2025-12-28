package taskappdeploy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	deploymentInfoCacheExp = 8 * time.Hour
)

type Executor struct {
	logger             logging.Logger
	settingRepo        repository.SettingRepo
	deploymentRepo     repository.DeploymentRepo
	deploymentInfoRepo cacherepository.DeploymentInfoRepo
	dockerManager      *docker.Manager
}

func NewExecutor(
	taskQueue taskqueue.TaskQueue,
	logger logging.Logger,
	settingRepo repository.SettingRepo,
	deploymentRepo repository.DeploymentRepo,
	deploymentInfoRepo cacherepository.DeploymentInfoRepo,
	dockerManager *docker.Manager,
) *Executor {
	p := &Executor{
		logger:             logger,
		settingRepo:        settingRepo,
		deploymentRepo:     deploymentRepo,
		deploymentInfoRepo: deploymentInfoRepo,
		dockerManager:      dockerManager,
	}
	taskQueue.RegisterExecutor(base.TaskTypeAppDeploy, p.execute)
	return p
}

type taskData struct {
	Task               *entity.Task
	Deployment         *entity.Deployment
	DeploymentCanceled bool
}

func (e *Executor) execute(
	ctx context.Context,
	db database.Tx,
	task *entity.Task,
) (err error) {
	data := &taskData{Task: task}
	deployment, err := e.loadDeployment(ctx, db, data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if deployment == nil {
		return nil
	}

	defer func() {
		if err == nil {
			if r := recover(); r != nil {
				err = apperrors.NewPanic(fmt.Sprintf("%v", r))
			}
		}
		_ = e.deploymentInfoRepo.Del(ctx, deployment.ID)
	}()

	var depErr error
	depSettings := deployment.DeploymentSettings
	switch {
	case depSettings.ImageSource != nil && depSettings.ImageSource.Enabled:
		depErr = e.deployFromImage(ctx, db, data)
	case depSettings.RepoSource != nil && depSettings.RepoSource.Enabled:
		depErr = e.deployFromRepo(ctx, db, data)
	case depSettings.TarballSource != nil && depSettings.TarballSource.Enabled:
		depErr = e.deployFromTarball(ctx, db, data)
	}

	if data.DeploymentCanceled {
		deployment.Status = base.DeploymentStatusCanceled
	} else {
		deployment.Status = gofn.If(depErr != nil, base.DeploymentStatusFailed, base.DeploymentStatusDone) //nolint
		deployment.EndedAt = timeutil.NowUTC()
	}

	err = e.updateDeployment(ctx, db, deployment)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (e *Executor) loadDeployment(
	ctx context.Context,
	db database.Tx,
	data *taskData,
) (*entity.Deployment, error) {
	task := data.Task
	args, err := task.ArgsAsAppDeploy()
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	deployment, err := e.deploymentRepo.GetByID(ctx, db, args.Deployment.ID,
		bunex.SelectWhereIn("deployment.status IN (?)",
			base.DeploymentStatusNotStarted, base.DeploymentStatusInProgress),
		bunex.SelectRelation("App",
			bunex.SelectColumns("id", "service_id"),
			bunex.SelectWhere("app.status = ?", base.AppStatusActive),
		),
		bunex.SelectFor("UPDATE OF deployment"),
	)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.Wrap(err)
	}
	if deployment == nil || deployment.App == nil { // no active deployment, return
		return nil, nil
	}

	if deployment.Status == base.DeploymentStatusNotStarted {
		deployment.StartedAt = data.Task.StartedAt
		deployment.Status = base.DeploymentStatusInProgress
	}

	// Put deployment status in redis
	err = e.deploymentInfoRepo.Set(ctx, deployment.ID, &cacheentity.DeploymentInfo{
		ID:        deployment.ID,
		AppID:     deployment.AppID,
		Status:    base.DeploymentStatusInProgress,
		StartedAt: deployment.StartedAt,
	}, deploymentInfoCacheExp)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	data.Deployment = deployment
	return deployment, nil
}

func (e *Executor) updateDeployment(
	ctx context.Context,
	db database.Tx,
	deployment *entity.Deployment,
) error {
	err := e.deploymentRepo.Update(ctx, db, deployment)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (e *Executor) isDeploymentCanceled(
	ctx context.Context,
	deployment *entity.Deployment,
) (bool, error) {
	info, err := e.deploymentInfoRepo.Get(ctx, deployment.ID)
	if err != nil {
		return false, apperrors.Wrap(err)
	}
	return info != nil && info.Cancel, nil
}

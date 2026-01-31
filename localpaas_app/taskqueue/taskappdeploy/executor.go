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
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/repository/cacherepository"
	"github.com/localpaas/localpaas/localpaas_app/service/envvarservice"
	"github.com/localpaas/localpaas/localpaas_app/taskqueue"
	"github.com/localpaas/localpaas/services/docker"
)

const (
	deploymentInfoCacheExp = 4 * time.Hour
)

type Executor struct {
	logger             logging.Logger
	redisClient        rediscache.Client
	settingRepo        repository.SettingRepo
	deploymentRepo     repository.DeploymentRepo
	deploymentLogRepo  repository.DeploymentLogRepo
	taskInfoRepo       cacherepository.TaskInfoRepo
	deploymentInfoRepo cacherepository.DeploymentInfoRepo
	dockerManager      *docker.Manager
	envVarService      envvarservice.EnvVarService
}

func NewExecutor(
	taskQueue taskqueue.TaskQueue,
	logger logging.Logger,
	redisClient rediscache.Client,
	settingRepo repository.SettingRepo,
	deploymentRepo repository.DeploymentRepo,
	deploymentLogRepo repository.DeploymentLogRepo,
	taskInfoRepo cacherepository.TaskInfoRepo,
	deploymentInfoRepo cacherepository.DeploymentInfoRepo,
	dockerManager *docker.Manager,
	envVarService envvarservice.EnvVarService,
) *Executor {
	p := &Executor{
		logger:             logger,
		redisClient:        redisClient,
		settingRepo:        settingRepo,
		deploymentRepo:     deploymentRepo,
		deploymentLogRepo:  deploymentLogRepo,
		taskInfoRepo:       taskInfoRepo,
		deploymentInfoRepo: deploymentInfoRepo,
		dockerManager:      dockerManager,
		envVarService:      envVarService,
	}
	taskQueue.RegisterExecutor(base.TaskTypeAppDeploy, p.execute)
	return p
}

type taskData struct {
	*taskqueue.TaskExecData
	Deployment       *entity.Deployment
	DeploymentOutput *entity.AppDeploymentOutput
	Step             string
	LogStore         *realtimelog.Store
}

func (e *Executor) execute(
	ctx context.Context,
	db database.Tx,
	task *taskqueue.TaskExecData,
) (err error) {
	data := &taskData{TaskExecData: task}
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
		_ = e.saveLogs(ctx, db, data)
	}()

	var depErr error
	depSettings := deployment.Settings
	switch depSettings.ActiveMethod {
	case base.DeploymentMethodImage:
		depErr = e.deployFromImage(ctx, db, data)
	case base.DeploymentMethodRepo:
		depErr = e.deployFromRepo(ctx, db, data)
	case base.DeploymentMethodTarball:
		depErr = e.deployFromTarball(ctx, db, data)
	}

	deployment.EndedAt = timeutil.NowUTC()
	if data.Canceled {
		deployment.Status = base.DeploymentStatusCanceled
	} else {
		deployment.Status = gofn.If(depErr != nil, base.DeploymentStatusFailed, base.DeploymentStatusDone)
		deployment.Output = data.DeploymentOutput
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

	deployment, err := e.deploymentRepo.GetByID(ctx, db, "", args.Deployment.ID,
		bunex.SelectWhereIn("deployment.status IN (?)",
			base.DeploymentStatusNotStarted, base.DeploymentStatusInProgress),
		bunex.SelectRelation("App",
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
		TaskID:    task.ID,
		Status:    base.DeploymentStatusInProgress,
		StartedAt: deployment.StartedAt,
	}, deploymentInfoCacheExp)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	data.Deployment = deployment
	data.DeploymentOutput = &entity.AppDeploymentOutput{}
	logStoreKey := fmt.Sprintf("deployment:%s:log", deployment.ID)
	data.LogStore = realtimelog.NewStore(logStoreKey, true, e.redisClient)

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

func (e *Executor) saveLogs(
	ctx context.Context,
	db database.Tx,
	taskData *taskData,
) error {
	deployment := taskData.Deployment
	logStore := taskData.LogStore
	if logStore == nil {
		return nil
	}

	duration := deployment.EndedAt.Sub(deployment.StartedAt)
	_ = logStore.Add(ctx, realtimelog.NewOutFrame("Deployment finished in "+duration.String(), nil))

	logFrames, err := logStore.GetData(ctx, 0)
	if err != nil {
		return apperrors.Wrap(err)
	}
	err = logStore.Close() //nolint
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Insert data in to DB by chunk to avoid exceeding DBMS limit
	for _, chunk := range gofn.Chunk(logFrames, 10000) { //nolint
		deploymentLogs := make([]*entity.DeploymentLog, 0, len(chunk))
		for _, logFrame := range chunk {
			deploymentLogs = append(deploymentLogs, &entity.DeploymentLog{
				DeploymentID: deployment.ID,
				Type:         logFrame.Type,
				Data:         logFrame.Data,
				Ts:           logFrame.Ts,
			})
		}
		err = e.deploymentLogRepo.InsertMulti(ctx, db, deploymentLogs)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (e *Executor) addStepStartLog(
	ctx context.Context,
	taskData *taskData,
	msg string,
) {
	_ = taskData.LogStore.Add(ctx,
		realtimelog.NewOutFrame("---------------------------------", nil),
		realtimelog.NewOutFrame(msg, nil))
}

func (e *Executor) addStepEndLog(
	ctx context.Context,
	taskData *taskData,
	start time.Time,
	err error,
) {
	duration := timeutil.NowUTC().Sub(start)
	if err != nil {
		_ = taskData.LogStore.Add(ctx, realtimelog.NewOutFrame("Task finished in "+duration.String()+
			" with error: "+err.Error(), nil))
	} else {
		_ = taskData.LogStore.Add(ctx, realtimelog.NewOutFrame("Task finished in "+duration.String(), nil))
	}
}

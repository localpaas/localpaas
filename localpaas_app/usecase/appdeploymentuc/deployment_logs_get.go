package appdeploymentuc

import (
	"context"
	"fmt"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc/appdeploymentdto"
)

const (
	deploymentLogBatchThresholdPeriod = time.Millisecond * 1000
	deploymentLogBatchMaxFrame        = 20
	deploymentLogSessionTimeout       = 10 * time.Minute
)

func (uc *AppDeploymentUC) GetDeploymentLogs(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdeploymentdto.GetDeploymentLogsReq,
) (*appdeploymentdto.GetDeploymentLogsResp, error) {
	deployment, err := uc.deploymentRepo.GetByID(ctx, uc.db, req.AppID, req.DeploymentID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if deployment.Status == base.DeploymentStatusNotStarted {
		return uc.getRealtimeDeploymentLogs(ctx, deployment, req)
	}

	return uc.getHistoryDeploymentLogs(ctx, deployment, req)
}

func (uc *AppDeploymentUC) getRealtimeDeploymentLogs(
	ctx context.Context,
	deployment *entity.Deployment,
	req *appdeploymentdto.GetDeploymentLogsReq,
) (*appdeploymentdto.GetDeploymentLogsResp, error) {
	key := fmt.Sprintf("deployment:%s:log", deployment.ID)
	consumer := realtimelog.NewConsumer(key, uc.redisClient)

	resp := &appdeploymentdto.DeploymentLogsDataResp{}
	var err error
	if req.Follow {
		// NOTE: we don't want to keep the log stream session forever
		ctx, _ = context.WithTimeout(ctx, deploymentLogSessionTimeout) //nolint:govet

		resp.LogChan, resp.LogChanCloser, err = consumer.Consume(ctx, batchrecvchan.Options{
			ThresholdPeriod: deploymentLogBatchThresholdPeriod,
			MaxItem:         deploymentLogBatchMaxFrame,
		})
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	} else {
		frames, err := consumer.GetAllData(ctx)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp.Logs = append(resp.Logs, frames...)
	}

	return &appdeploymentdto.GetDeploymentLogsResp{
		Data: resp,
	}, nil
}

func (uc *AppDeploymentUC) getHistoryDeploymentLogs(
	ctx context.Context,
	deployment *entity.Deployment,
	req *appdeploymentdto.GetDeploymentLogsReq,
) (*appdeploymentdto.GetDeploymentLogsResp, error) {
	var listOpts []bunex.SelectQueryOption

	reverseLogs := false
	if req.Tail > 0 {
		listOpts = append(listOpts, bunex.SelectLimit(req.Tail),
			bunex.SelectOrder("id DESC"))
		reverseLogs = true
	} else {
		listOpts = append(listOpts, bunex.SelectOrder("id"))
	}

	if !req.Since.IsZero() {
		listOpts = append(listOpts,
			bunex.SelectWhere("task_log.ts >= ?", req.Since))
		if req.Duration > 0 {
			listOpts = append(listOpts,
				bunex.SelectWhere("task_log.ts < ?", req.Since.Add(req.Duration)))
		}
	}

	logs, _, err := uc.taskLogRepo.List(ctx, uc.db, "", deployment.ID, nil, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	// Reverse the data
	if reverseLogs {
		gofn.Reverse(logs)
	}

	logFrames := appdeploymentdto.TransformDeploymentLogs(logs)
	logChan := make(chan []*realtimelog.LogFrame, 100) //nolint:mnd

	resp := &appdeploymentdto.DeploymentLogsDataResp{
		Logs:          logFrames,
		LogChan:       logChan,
		LogChanCloser: func() error { return nil },
	}

	go func() {
		for _, chunk := range gofn.Chunk(logFrames, deploymentLogBatchMaxFrame) {
			logChan <- chunk
		}
		for len(logChan) > 0 {
			time.Sleep(300 * time.Millisecond) //nolint:mnd
		}
		close(logChan)
	}()

	return &appdeploymentdto.GetDeploymentLogsResp{
		Data: resp,
	}, nil
}

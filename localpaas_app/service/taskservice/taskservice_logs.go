package taskservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/batchrecvchan"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc/taskdto"
)

type GetTaskLogsReq struct {
	TaskID     string
	Follow     bool
	Since      time.Time
	Duration   time.Duration
	Tail       int
	Timestamps bool

	LogBatchThresholdPeriod time.Duration
	LogBatchMaxFrame        int
	LogSessionTimeout       time.Duration
}

type GetTaskLogsResp struct {
	Logs          []*applog.LogFrame
	LogChan       <-chan []*applog.LogFrame
	LogChanCloser func() error
}

func (s *taskService) GetTaskLogs(
	ctx context.Context,
	db database.IDB,
	req *GetTaskLogsReq,
) (*GetTaskLogsResp, error) {
	task, err := s.taskRepo.GetByID(ctx, db, "", req.TaskID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp := &GetTaskLogsResp{}
	err = s.queryLogsInDB(ctx, db, task, req, resp)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	err = s.queryRealtimeLogs(ctx, task, req, resp)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

func (s *taskService) queryRealtimeLogs(
	ctx context.Context,
	task *entity.Task,
	req *GetTaskLogsReq,
	resp *GetTaskLogsResp,
) error {
	if task.IsDone() || task.IsCanceled() || task.IsFailedCompletely() {
		return nil
	}

	taskInfo, err := s.taskInfoRepo.Get(ctx, task.ID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return apperrors.Wrap(err)
	}
	if taskInfo == nil {
		return nil
	}

	key := fmt.Sprintf("task:%s:log", task.ID)
	consumer := applog.NewConsumer(key, s.redisClient)

	if req.Follow {
		// NOTE: we don't want to keep the log stream session forever
		ctx, _ = context.WithTimeout(ctx, req.LogSessionTimeout) //nolint:govet

		resp.LogChan, resp.LogChanCloser, err = consumer.StartConsuming(ctx, batchrecvchan.Options{
			ThresholdPeriod: req.LogBatchThresholdPeriod,
			MaxItem:         req.LogBatchMaxFrame,
		})
		if err != nil {
			return apperrors.Wrap(err)
		}
	} else {
		frames, err := consumer.GetAllData(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
		resp.Logs = append(resp.Logs, frames...)
	}
	return nil
}

func (s *taskService) queryLogsInDB(
	ctx context.Context,
	db database.IDB,
	task *entity.Task,
	req *GetTaskLogsReq,
	resp *GetTaskLogsResp,
) error {
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

	logs, _, err := s.taskLogRepo.List(ctx, db, task.ID, "", nil, listOpts...)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Reverse the data
	if reverseLogs {
		gofn.Reverse(logs)
	}

	resp.Logs = append(resp.Logs, taskdto.TransformTaskLogs(logs)...)
	return nil
}

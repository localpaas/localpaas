package taskuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc/taskdto"
)

func (uc *TaskUC) GetTask(
	ctx context.Context,
	auth *basedto.Auth,
	req *taskdto.GetTaskReq,
) (*taskdto.GetTaskResp, error) {
	task, err := uc.taskRepo.GetByID(ctx, uc.db, "", req.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	var taskInfo *cacheentity.TaskInfo
	if task.Status != base.TaskStatusDone && task.Status != base.TaskStatusCanceled {
		taskInfo, err = uc.cacheTaskInfoRepo.Get(ctx, task.ID)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	resp, err := taskdto.TransformTask(task, taskInfo)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &taskdto.GetTaskResp{
		Data: resp,
	}, nil
}

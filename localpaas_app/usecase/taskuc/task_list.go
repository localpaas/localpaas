package taskuc

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc/taskdto"
)

func (uc *TaskUC) ListTask(
	ctx context.Context,
	auth *basedto.Auth,
	req *taskdto.ListTaskReq,
) (*taskdto.ListTaskResp, error) {
	taskInfoMap, err := uc.cacheTaskInfoRepo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	inprogressTaskIDs := make([]string, 0, len(taskInfoMap))
	for id, info := range taskInfoMap {
		if info.Status == base.TaskStatusInProgress {
			inprogressTaskIDs = append(inprogressTaskIDs, id)
		}
	}

	var listOpts []bunex.SelectQueryOption
	if len(req.JobID) > 0 {
		listOpts = append(listOpts, bunex.SelectWhereIn("task.job_id IN (?)", req.JobID...))
	}
	if len(req.Status) > 0 { //nolint:nestif
		statuses := req.Status
		if gofn.Contain(statuses, base.TaskStatusInProgress) {
			cond := bunex.SelectWhereIn("task.id IN (?)", inprogressTaskIDs...)
			statuses = gofn.Drop(statuses, base.TaskStatusInProgress)
			if len(statuses) == 0 {
				listOpts = append(listOpts, cond)
			} else {
				listOpts = append(listOpts, cond,
					bunex.SelectWhereOrGroup(
						bunex.SelectWhereNotIn("task.id NOT IN (?)", inprogressTaskIDs...),
						bunex.SelectWhereIn("task.status IN (?)", statuses),
					),
				)
			}
		} else {
			listOpts = append(listOpts,
				bunex.SelectWhereNotIn("task.id NOT IN (?)", inprogressTaskIDs...),
				bunex.SelectWhereIn("task.status IN (?)", statuses...))
		}
	}
	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("task.type ILIKE ?", keyword),
			),
		)
	}
	if len(auth.AllowObjectIDs) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhereIn("task.id IN (?)", auth.AllowObjectIDs...),
		)
	}

	tasks, paging, err := uc.taskRepo.List(ctx, uc.db, "", &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := taskdto.TransformTasks(tasks, taskInfoMap)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &taskdto.ListTaskResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}

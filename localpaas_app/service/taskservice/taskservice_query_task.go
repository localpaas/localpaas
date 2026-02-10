package taskservice

import (
	"context"
	"errors"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

type GetTaskReq struct {
	ID       string
	Type     base.TaskType
	JobID    string
	TargetID string

	SkipQueryCache bool
}

type GetTaskResp struct {
	Task     *entity.Task
	TaskInfo *cacheentity.TaskInfo
}

func (s *taskService) GetTask(
	ctx context.Context,
	db database.IDB,
	req *GetTaskReq,
	extraOpts ...bunex.SelectQueryOption,
) (*GetTaskResp, error) {
	var getOpts []bunex.SelectQueryOption
	if req.JobID != "" {
		getOpts = append(getOpts, bunex.SelectWhere("task.job_id = ?", req.JobID))
	}
	if req.TargetID != "" {
		getOpts = append(getOpts, bunex.SelectWhere("task.target_id = ?", req.TargetID))
	}
	getOpts = append(getOpts, extraOpts...)

	task, err := s.taskRepo.GetByID(ctx, db, req.Type, req.ID, getOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	var taskInfo *cacheentity.TaskInfo
	if !req.SkipQueryCache && !task.IsDone() && !task.IsCanceled() && !task.IsFailedCompletely() {
		taskInfo, err = s.taskInfoRepo.Get(ctx, task.ID)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.Wrap(err)
		}
	}

	return &GetTaskResp{
		Task:     task,
		TaskInfo: taskInfo,
	}, nil
}

type ListTaskReq struct {
	JobID  []string
	Status []base.TaskStatus
	Search string
	Paging basedto.Paging

	SkipQueryCache bool
}

type ListTaskResp struct {
	PagingMeta  *basedto.PagingMeta
	Tasks       []*entity.Task
	TaskInfoMap map[string]*cacheentity.TaskInfo
}

func (s *taskService) ListTask(
	ctx context.Context,
	db database.IDB,
	req *ListTaskReq,
	extraOpts ...bunex.SelectQueryOption,
) (_ *ListTaskResp, err error) {
	var taskInfoMap map[string]*cacheentity.TaskInfo
	var inprogressTaskIDs []string
	if !req.SkipQueryCache {
		taskInfoMap, err = s.taskInfoRepo.GetAll(ctx)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		inprogressTaskIDs = make([]string, 0, len(taskInfoMap))
		for id, info := range taskInfoMap {
			if info.Status == base.TaskStatusInProgress {
				inprogressTaskIDs = append(inprogressTaskIDs, id)
			}
		}
	}

	var listOpts []bunex.SelectQueryOption
	if len(req.JobID) > 0 {
		listOpts = append(listOpts, bunex.SelectWhereIn("task.job_id IN (?)", req.JobID...))
	}
	if len(req.Status) > 0 { //nolint:nestif
		statuses := req.Status
		if gofn.Contain(statuses, base.TaskStatusInProgress) {
			statuses = gofn.Drop(statuses, base.TaskStatusInProgress)
			if len(statuses) == 0 {
				listOpts = append(listOpts,
					bunex.SelectWhereIn("task.id IN (?)", inprogressTaskIDs...),
				)
			} else {
				listOpts = append(listOpts,
					bunex.SelectWhereGroup(
						bunex.SelectWhereIn("task.id IN (?)", inprogressTaskIDs...),
						bunex.SelectWhereOrGroup(
							bunex.SelectWhereNotIn("task.id NOT IN (?)", inprogressTaskIDs...),
							bunex.SelectWhereIn("task.status IN (?)", statuses),
						),
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
	listOpts = append(listOpts, extraOpts...)

	tasks, paging, err := s.taskRepo.List(ctx, db, "", &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ListTaskResp{
		PagingMeta:  paging,
		Tasks:       tasks,
		TaskInfoMap: taskInfoMap,
	}, nil
}

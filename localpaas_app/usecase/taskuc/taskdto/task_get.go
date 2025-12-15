package taskdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetTaskReq struct {
	ID string `json:"-"`
}

func NewGetTaskReq() *GetTaskReq {
	return &GetTaskReq{}
}

func (req *GetTaskReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetTaskResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *TaskResp         `json:"data"`
}

type TaskResp struct {
	ID             string            `json:"id"`
	Type           base.TaskType     `json:"type"`
	Status         base.TaskStatus   `json:"status"`
	Priority       base.TaskPriority `json:"priority"`
	MaxRetry       int               `json:"maxRetry"`
	RetryDelaySecs int               `json:"retryDelaySecs"`
	Command        string            `json:"command"`
	UpdateVer      int               `json:"updateVer"`

	RunAt     time.Time `json:"runAt"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func TransformTask(task *entity.Task, taskInfo *cacheentity.TaskInfo) (resp *TaskResp, err error) {
	if err = copier.Copy(&resp, &task); err != nil {
		return nil, apperrors.Wrap(err)
	}
	if taskInfo != nil {
		resp.Status = taskInfo.Status
		if taskInfo.Status == base.TaskStatusInProgress {
			resp.StartedAt = taskInfo.StartedAt
		}
	}
	return resp, nil
}

package systemcleanupdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/taskuc/taskdto"
)

type GetSystemCleanupTaskReq struct {
	settings.BaseSettingReq
	JobID  string `json:"-"`
	TaskID string `json:"-"`
}

func NewGetSystemCleanupTaskReq() *GetSystemCleanupTaskReq {
	return &GetSystemCleanupTaskReq{}
}

func (req *GetSystemCleanupTaskReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.JobID, true, "jobId")...)
	validators = append(validators, basedto.ValidateID(&req.TaskID, true, "taskId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSystemCleanupTaskResp struct {
	Meta *basedto.Meta     `json:"meta"`
	Data *taskdto.TaskResp `json:"data"`
}

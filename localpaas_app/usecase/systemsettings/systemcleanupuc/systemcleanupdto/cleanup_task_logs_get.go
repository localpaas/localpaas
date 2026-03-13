package systemcleanupdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetSystemCleanupTaskLogsReq struct {
	settings.BaseSettingReq
	JobID    string        `json:"-"`
	TaskID   string        `json:"-"`
	Follow   bool          `json:"-" mapstructure:"follow"`
	Since    time.Time     `json:"-" mapstructure:"since"`
	Duration time.Duration `json:"-" mapstructure:"duration"`
	Tail     int           `json:"-" mapstructure:"tail"`
}

func NewGetSystemCleanupTaskLogsReq() *GetSystemCleanupTaskLogsReq {
	return &GetSystemCleanupTaskLogsReq{}
}

func (req *GetSystemCleanupTaskLogsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.JobID, true, "jobId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSystemCleanupTaskLogsResp struct {
	Meta *basedto.Meta                  `json:"meta"`
	Data *SystemCleanupTaskLogsDataResp `json:"data"`
}

type SystemCleanupTaskLogsDataResp struct {
	Logs          []*applog.LogFrame        `json:"logs"`
	LogChan       <-chan []*applog.LogFrame `json:"-"`
	LogChanCloser func() error              `json:"-"`
}

func TransformSystemCleanupTaskLogs(logs []*entity.TaskLog) (resp []*applog.LogFrame) {
	resp = make([]*applog.LogFrame, 0, len(logs))
	for _, log := range logs {
		resp = append(resp, &applog.LogFrame{
			Type: log.Type,
			Data: log.Data,
			Ts:   log.Ts,
		})
	}
	return resp
}

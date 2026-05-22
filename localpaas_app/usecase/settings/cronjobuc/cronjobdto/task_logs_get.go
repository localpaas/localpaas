package cronjobdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tasklog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetCronJobTaskLogsReq struct {
	settings.BaseSettingReq
	JobID    string            `json:"-"`
	TaskID   string            `json:"-"`
	Follow   bool              `json:"-" mapstructure:"follow"`
	Since    time.Time         `json:"-" mapstructure:"since"`
	Duration timeutil.Duration `json:"-" mapstructure:"duration"`
	Tail     int               `json:"-" mapstructure:"tail"`
}

func NewGetCronJobTaskLogsReq() *GetCronJobTaskLogsReq {
	return &GetCronJobTaskLogsReq{}
}

func (req *GetCronJobTaskLogsReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.JobID, true, "jobId")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetCronJobTaskLogsResp struct {
	Meta *basedto.Meta            `json:"meta"`
	Data *CronJobTaskLogsDataResp `json:"data"`
}

type CronJobTaskLogsDataResp struct {
	Logs          []*tasklog.LogFrame        `json:"logs"`
	LogChan       <-chan []*tasklog.LogFrame `json:"-"`
	LogChanCloser func() error               `json:"-"`
}

func TransformCronJobTaskLogs(logs []*entity.TaskLog) (resp []*tasklog.LogFrame) {
	resp = make([]*tasklog.LogFrame, 0, len(logs))
	for _, log := range logs {
		resp = append(resp, &tasklog.LogFrame{
			Type: log.Type,
			Data: log.Data,
			Ts:   log.Ts,
		})
	}
	return resp
}

package cronjobdto

import (
	"strings"

	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateCronJobReq struct {
	settings.CreateSettingReq
	*CronJobBaseReq
}

type CronJobBaseReq struct {
	Name       string                      `json:"name"`
	Kind       base.TaskType               `json:"kind"`
	CronType   base.CronJobType            `json:"cronType"`
	CronExpr   string                      `json:"cronExpr"`
	App        basedto.ObjectIDReq         `json:"app"`
	Priority   base.TaskPriority           `json:"priority"`
	MaxRetry   int                         `json:"maxRetry"`
	RetryDelay timeutil.Duration           `json:"retryDelay"`
	Timeout    timeutil.Duration           `json:"timeout"`
	Command    *CronJobContainerCommandReq `json:"command"`
}

type CronJobContainerCommandReq struct {
	Command    string `json:"command"`
	WorkingDir string `json:"workingDir"`
}

func (req *CronJobBaseReq) ToEntity() *entity.CronJob {
	item := &entity.CronJob{
		CronType:    req.CronType,
		CronExpr:    req.CronExpr,
		App:         entity.ObjectID{ID: req.App.ID},
		InitialTime: timeutil.NowUTC(),
		Priority:    req.Priority,
		MaxRetry:    req.MaxRetry,
		RetryDelay:  req.RetryDelay,
		Timeout:     req.Timeout,
	}
	if req.Command != nil {
		item.Command = &entity.CronJobContainerCommand{
			Command:    req.Command.Command,
			WorkingDir: req.Command.WorkingDir,
		}
	}
	return item
}

func (req *CronJobBaseReq) modifyRequest() error {
	req.Name = strings.TrimSpace(req.Name)
	req.CronExpr = strings.TrimSpace(req.CronExpr)
	req.Priority = gofn.Coalesce(req.Priority, base.TaskPriorityDefault)
	return nil
}

func (req *CronJobBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateCronJobReq() *CreateCronJobReq {
	return &CreateCronJobReq{}
}

func (req *CreateCronJobReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *CreateCronJobReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateCronJobResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

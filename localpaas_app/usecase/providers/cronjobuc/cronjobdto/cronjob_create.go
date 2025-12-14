package cronjobdto

import (
	"strings"

	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateCronJobReq struct {
	*CronJobBaseReq
}

type CronJobBaseReq struct {
	Name           string            `json:"name"`
	Kind           base.TaskType     `json:"kind"`
	Cron           string            `json:"cron"`
	Priority       base.TaskPriority `json:"priority"`
	MaxRetry       int               `json:"maxRetry"`
	RetryDelaySecs int               `json:"retryDelaySecs"`
	Command        string            `json:"command"`
}

func (req *CronJobBaseReq) modifyRequest() error {
	req.Name = strings.TrimSpace(req.Name)
	req.Cron = strings.TrimSpace(req.Cron)
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
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

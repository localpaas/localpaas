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
	"github.com/localpaas/localpaas/localpaas_app/usecase/notification/notificationdto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type CreateCronJobReq struct {
	settings.CreateSettingReq
	*CronJobBaseReq
}

type CronJobBaseReq struct {
	Name         string                                       `json:"name"`
	CronType     base.CronJobType                             `json:"cronType"`
	CronExpr     string                                       `json:"cronExpr"`
	App          basedto.ObjectIDReq                          `json:"app"`
	Priority     base.TaskPriority                            `json:"priority"`
	MaxRetry     int                                          `json:"maxRetry"`
	RetryDelay   timeutil.Duration                            `json:"retryDelay"`
	Timeout      timeutil.Duration                            `json:"timeout"`
	Command      *CronJobContainerCommandReq                  `json:"command"`
	Notification *notificationdto.DefaultResultNtfnSettingReq `json:"notification"`
}

func (req *CronJobBaseReq) ToEntity() *entity.CronJob {
	return &entity.CronJob{
		CronType:     req.CronType,
		CronExpr:     req.CronExpr,
		App:          entity.ObjectID{ID: req.App.ID},
		InitialTime:  timeutil.NowUTC(),
		Priority:     req.Priority,
		MaxRetry:     req.MaxRetry,
		RetryDelay:   req.RetryDelay,
		Timeout:      req.Timeout,
		Command:      req.Command.ToEntity(),
		Notification: req.Notification.ToEntity(),
	}
}

type CronJobContainerCommandReq struct {
	RunInShell string                       `json:"runInShell"`
	Command    string                       `json:"command"`
	WorkingDir string                       `json:"workingDir"`
	EnvVars    []*basedto.EnvVarReq         `json:"envVars"`
	ArgGroups  []*CronJobCommandArgGroupReq `json:"argGroups"`
}

func (req *CronJobContainerCommandReq) ToEntity() *entity.CronJobContainerCommand {
	if req == nil {
		return nil
	}
	return &entity.CronJobContainerCommand{
		Command:    req.Command,
		WorkingDir: req.WorkingDir,
		ArgGroups: gofn.MapSlice(req.ArgGroups, func(item *CronJobCommandArgGroupReq) *entity.CronJobCommandArgGroup {
			return item.ToEntity()
		}),
	}
}

type CronJobCommandArgGroupReq struct {
	ExportEnv string                  `json:"exportEnv"`
	Separator string                  `json:"separator"`
	Args      []*CronJobCommandArgReq `json:"args,omitempty"`
}

func (req *CronJobCommandArgGroupReq) ToEntity() *entity.CronJobCommandArgGroup {
	if req == nil {
		return nil
	}
	return &entity.CronJobCommandArgGroup{
		ExportEnv: req.ExportEnv,
		Separator: req.Separator,
		Args: gofn.MapSlice(req.Args, func(item *CronJobCommandArgReq) *entity.CronJobCommandArg {
			return item.ToEntity()
		}),
	}
}

type CronJobCommandArgReq struct {
	Use   bool   `json:"use"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (req *CronJobCommandArgReq) ToEntity() *entity.CronJobCommandArg {
	if req == nil {
		return nil
	}
	return &entity.CronJobCommandArg{
		Use:   req.Use,
		Name:  req.Name,
		Value: req.Value,
	}
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

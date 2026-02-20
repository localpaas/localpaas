package cronjobdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetCronJobReq struct {
	settings.GetSettingReq
}

func NewGetCronJobReq() *GetCronJobReq {
	return &GetCronJobReq{}
}

func (req *GetCronJobReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetCronJobResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data *CronJobResp  `json:"data"`
}

type CronJobResp struct {
	*settings.BaseSettingResp
	CronType     base.CronJobType             `json:"cronType"`
	CronExpr     string                       `json:"cronExpr"`
	App          *basedto.NamedObjectResp     `json:"app"`
	InitialTime  time.Time                    `json:"initialTime"`
	Priority     base.TaskPriority            `json:"priority"`
	MaxRetry     int                          `json:"maxRetry"`
	RetryDelay   timeutil.Duration            `json:"retryDelay"`
	Timeout      timeutil.Duration            `json:"timeout"`
	Command      *CronJobContainerCommandResp `json:"command"`
	Notification *CronJobNotificationResp     `json:"notification"`
}

type CronJobContainerCommandResp struct {
	RunInShell string                        `json:"runInShell"`
	Command    string                        `json:"command"`
	WorkingDir string                        `json:"workingDir"`
	EnvVars    []*basedto.EnvVarResp         `json:"envVars"`
	ArgGroups  []*CronJobCommandArgGroupResp `json:"argGroups"`
}

type CronJobCommandArgGroupResp struct {
	ExportEnv string                   `json:"exportEnv"`
	Separator string                   `json:"separator"`
	Args      []*CronJobCommandArgResp `json:"args"`
}

type CronJobCommandArgResp struct {
	Use   bool   `json:"use"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CronJobNotificationResp struct {
	Success *settings.BaseSettingResp `json:"success"`
	Failure *settings.BaseSettingResp `json:"failure"`
}

func TransformCronJob(
	setting *entity.Setting,
	refObjects *entity.RefObjects,
) (resp *CronJobResp, err error) {
	config := setting.MustAsCronJob()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	if config.App.ID != "" {
		refApp := refObjects.RefApps[config.App.ID]
		if err = copier.Copy(&resp.App, refApp); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	if resp.Notification != nil {
		if resp.Notification.Success != nil {
			itemResp, _ := settings.TransformSettingBase(refObjects.RefSettings[resp.Notification.Success.ID])
			resp.Notification.Success = itemResp
		}
		if resp.Notification.Failure != nil {
			itemResp, _ := settings.TransformSettingBase(refObjects.RefSettings[resp.Notification.Failure.ID])
			resp.Notification.Failure = itemResp
		}
	}

	return resp, nil
}

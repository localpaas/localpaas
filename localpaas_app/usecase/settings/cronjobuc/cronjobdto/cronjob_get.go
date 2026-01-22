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
	Meta *basedto.BaseMeta `json:"meta"`
	Data *CronJobResp      `json:"data"`
}

type CronJobResp struct {
	*settings.BaseSettingResp
	Cron        string            `json:"cron"`
	InitialTime time.Time         `json:"initialTime"`
	Priority    base.TaskPriority `json:"priority"`
	MaxRetry    int               `json:"maxRetry"`
	RetryDelay  timeutil.Duration `json:"retryDelay"`
	Timeout     timeutil.Duration `json:"timeout"`
	Command     string            `json:"command"`
}

func TransformCronJob(setting *entity.Setting, objectID string) (resp *CronJobResp, err error) {
	config := setting.MustAsCronJob()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting, objectID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return resp, nil
}

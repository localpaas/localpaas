package cronjobdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
)

type GetCronJobReq struct {
	ID string `json:"-"`
}

func NewGetCronJobReq() *GetCronJobReq {
	return &GetCronJobReq{}
}

func (req *GetCronJobReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ID, true, "id")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetCronJobResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
	Data *CronJobResp      `json:"data"`
}

type CronJobResp struct {
	ID             string             `json:"id"`
	Kind           string             `json:"kind"`
	Name           string             `json:"name"`
	Status         base.SettingStatus `json:"status"`
	Cron           string             `json:"cron"`
	InitialTime    time.Time          `json:"initialTime"`
	Priority       base.TaskPriority  `json:"priority"`
	MaxRetry       int                `json:"maxRetry"`
	RetryDelaySecs int                `json:"retryDelaySecs"`
	Command        string             `json:"command"`
	UpdateVer      int                `json:"updateVer"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	ExpireAt  *time.Time `json:"expireAt,omitempty" copy:",nilonzero"`
}

func TransformCronJob(setting *entity.Setting) (resp *CronJobResp, err error) {
	if err = copier.Copy(&resp, &setting); err != nil {
		return nil, apperrors.Wrap(err)
	}

	config := setting.MustAsCronJob()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	return resp, nil
}

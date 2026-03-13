package systemcleanupdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/copier"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type GetSystemCleanupReq struct {
	settings.GetSettingReq
}

func NewGetSystemCleanupReq() *GetSystemCleanupReq {
	return &GetSystemCleanupReq{}
}

func (req *GetSystemCleanupReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type GetSystemCleanupResp struct {
	Meta *basedto.Meta      `json:"meta"`
	Data *SystemCleanupResp `json:"data"`
}

type SystemCleanupResp struct {
	*settings.BaseSettingResp
	ScheduleInterval  timeutil.Duration      `json:"scheduleInterval"`
	ScheduleFrom      time.Time              `json:"scheduleFrom"`
	DBObjectRetention *DBObjectRetentionResp `json:"dbObjectRetention"`
	ClusterCleanup    *ClusterCleanupResp    `json:"clusterCleanup"`
	Notification      *NotificationResp      `json:"notification"`
}

type DBObjectRetentionResp struct {
	Enabled        bool              `json:"enabled"`
	Tasks          timeutil.Duration `json:"tasks"`
	SysErrors      timeutil.Duration `json:"sysErrors"`
	Deployments    timeutil.Duration `json:"deployments"`
	DeletedObjects timeutil.Duration `json:"deletedObjects"`
}

type ClusterCleanupResp struct {
	Enabled         bool `json:"enabled"`
	PruneImages     bool `json:"pruneImages"`
	PruneVolumes    bool `json:"pruneVolumes"`
	PruneNetworks   bool `json:"pruneNetworks"`
	PruneContainers bool `json:"pruneContainers"`
}

type NotificationResp struct {
	Success *settings.BaseSettingResp `json:"success"`
	Failure *settings.BaseSettingResp `json:"failure"`
}

func TransformSystemCleanup(
	setting *entity.Setting,
	refObjects *entity.RefObjects,
) (resp *SystemCleanupResp, err error) {
	config := setting.MustAsSystemCleanup()
	if err = copier.Copy(&resp, config); err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.BaseSettingResp, err = settings.TransformSettingBase(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
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

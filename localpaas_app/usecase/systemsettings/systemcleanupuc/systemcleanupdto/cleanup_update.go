package systemcleanupdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type UpdateSystemCleanupReq struct {
	settings.UpdateSettingReq
	*SystemCleanupBaseReq
}

type SystemCleanupBaseReq struct {
	ScheduleInterval  timeutil.Duration     `json:"scheduleInterval"`
	ScheduleFrom      time.Time             `json:"scheduleFrom"`
	DBObjectRetention *DBObjectRetentionReq `json:"dbObjectRetention"`
	ClusterCleanup    *ClusterCleanupReq    `json:"clusterCleanup"`
	Notification      *NotificationReq      `json:"notification"`
}

func (req *SystemCleanupBaseReq) ToEntity() *entity.SystemCleanup {
	return &entity.SystemCleanup{
		ScheduleInterval:  req.ScheduleInterval,
		ScheduleFrom:      req.ScheduleFrom,
		DBObjectRetention: req.DBObjectRetention.ToEntity(),
		ClusterCleanup:    req.ClusterCleanup.ToEntity(),
		Notification:      req.Notification.ToEntity(),
	}
}

type DBObjectRetentionReq struct {
	Enabled        bool              `json:"enabled"`
	Tasks          timeutil.Duration `json:"tasks"`
	SysErrors      timeutil.Duration `json:"sysErrors"`
	Deployments    timeutil.Duration `json:"deployments"`
	DeletedObjects timeutil.Duration `json:"deletedObjects"`
}

func (req *DBObjectRetentionReq) ToEntity() *entity.DBObjectRetention {
	if req == nil {
		return nil
	}
	return &entity.DBObjectRetention{
		Enabled:        req.Enabled,
		Tasks:          req.Tasks,
		SysErrors:      req.SysErrors,
		Deployments:    req.Deployments,
		DeletedObjects: req.DeletedObjects,
	}
}

type ClusterCleanupReq struct {
	Enabled         bool `json:"enabled"`
	PruneImages     bool `json:"pruneImages"`
	PruneVolumes    bool `json:"pruneVolumes"`
	PruneNetworks   bool `json:"pruneNetworks"`
	PruneContainers bool `json:"pruneContainers"`
}

func (req *ClusterCleanupReq) ToEntity() *entity.ClusterCleanup {
	if req == nil {
		return nil
	}
	return &entity.ClusterCleanup{
		Enabled:         req.Enabled,
		PruneImages:     req.PruneImages,
		PruneVolumes:    req.PruneVolumes,
		PruneNetworks:   req.PruneNetworks,
		PruneContainers: req.PruneContainers,
	}
}

type NotificationReq struct {
	Success basedto.ObjectIDReq `json:"success"`
	Failure basedto.ObjectIDReq `json:"failure"`
}

func (req *NotificationReq) ToEntity() *entity.CronJobNotification {
	if req == nil {
		return nil
	}
	return &entity.CronJobNotification{
		Success: entity.ObjectID{ID: req.Success.ID},
		Failure: entity.ObjectID{ID: req.Failure.ID},
	}
}

func (req *SystemCleanupBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewUpdateSystemCleanupReq() *UpdateSystemCleanupReq {
	return &UpdateSystemCleanupReq{}
}

func (req *UpdateSystemCleanupReq) ModifyRequest() error {
	if !req.ScheduleFrom.IsZero() {
		req.ScheduleFrom = req.ScheduleFrom.Truncate(time.Minute)
	}
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSystemCleanupReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSystemCleanupResp struct {
	Meta *basedto.Meta `json:"meta"`
}

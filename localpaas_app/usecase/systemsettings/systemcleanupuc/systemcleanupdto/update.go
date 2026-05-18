package systemcleanupdto

import (
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
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
	Status            base.SettingStatus                `json:"status"`
	ScheduleInterval  timeutil.Duration                 `json:"scheduleInterval"`
	ScheduleFrom      time.Time                         `json:"scheduleFrom"`
	DBObjectRetention DBObjectRetentionReq              `json:"dbObjectRetention"`
	ClusterCleanup    SystemClusterCleanupReq           `json:"clusterCleanup"`
	BackupCleanup     SystemBackupCleanupReq            `json:"backupCleanup"`
	Notification      *basedto.BaseEventNotificationReq `json:"notification"`
}

func (req *SystemCleanupBaseReq) ToEntity() *entity.SystemCleanup {
	return &entity.SystemCleanup{
		ScheduleInterval:  req.ScheduleInterval,
		ScheduleFrom:      req.ScheduleFrom,
		DBObjectRetention: req.DBObjectRetention.ToEntity(),
		ClusterCleanup:    req.ClusterCleanup.ToEntity(),
		BackupCleanup:     req.BackupCleanup.ToEntity(),
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

func (req *DBObjectRetentionReq) ToEntity() entity.DBObjectRetention {
	return entity.DBObjectRetention{
		Enabled:        req.Enabled,
		Tasks:          req.Tasks,
		SysErrors:      req.SysErrors,
		Deployments:    req.Deployments,
		DeletedObjects: req.DeletedObjects,
	}
}

type SystemClusterCleanupReq struct {
	Enabled         bool `json:"enabled"`
	PruneImages     bool `json:"pruneImages"`
	PruneVolumes    bool `json:"pruneVolumes"`
	PruneNetworks   bool `json:"pruneNetworks"`
	PruneContainers bool `json:"pruneContainers"`
}

func (req *SystemClusterCleanupReq) ToEntity() entity.SystemClusterCleanup {
	return entity.SystemClusterCleanup{
		Enabled:         req.Enabled,
		PruneImages:     req.PruneImages,
		PruneVolumes:    req.PruneVolumes,
		PruneNetworks:   req.PruneNetworks,
		PruneContainers: req.PruneContainers,
	}
}

type SystemBackupCleanupReq struct {
	Enabled              bool              `json:"enabled"`
	CloudBackupRetention timeutil.Duration `json:"cloudBackupRetention"`
	LocalBackupRetention timeutil.Duration `json:"localBackupRetention"`
}

func (req *SystemBackupCleanupReq) ToEntity() entity.SystemBackupCleanup {
	return entity.SystemBackupCleanup{
		Enabled:              req.Enabled,
		CloudBackupRetention: req.CloudBackupRetention,
		LocalBackupRetention: req.LocalBackupRetention,
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

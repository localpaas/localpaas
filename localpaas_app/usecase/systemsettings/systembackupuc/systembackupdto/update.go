package systembackupdto

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

type UpdateSystemBackupReq struct {
	settings.UpdateSettingReq
	*SystemBackupBaseReq
}

type SystemBackupBaseReq struct {
	Status           base.SettingStatus                `json:"status"`
	ScheduleInterval timeutil.Duration                 `json:"scheduleInterval"`
	ScheduleFrom     time.Time                         `json:"scheduleFrom"`
	Compression      SystemBackupCompressionReq        `json:"compression"`
	Encryption       SystemBackupEncryptionReq         `json:"encryption"`
	CloudStorage     SystemBackupCloudStorageReq       `json:"cloudStorage"`
	DBBackupConfig   SystemBackupDBConfigReq           `json:"dbBackupConfig"`
	Notification     *basedto.BaseEventNotificationReq `json:"notification"`
}

func (req *SystemBackupBaseReq) ToEntity() *entity.SystemBackup {
	return &entity.SystemBackup{
		ScheduleInterval: req.ScheduleInterval,
		ScheduleFrom:     req.ScheduleFrom,
		Compression:      req.Compression.ToEntity(),
		Encryption:       req.Encryption.ToEntity(),
		CloudStorage:     req.CloudStorage.ToEntity(),
		DBBackupConfig:   req.DBBackupConfig.ToEntity(),
		Notification:     req.Notification.ToEntity(),
	}
}

type SystemBackupCompressionReq struct {
	Format base.FileCompressionFormat `json:"format,omitempty"`
}

func (req *SystemBackupCompressionReq) ToEntity() entity.SystemBackupCompression {
	return entity.SystemBackupCompression{
		Format: req.Format,
	}
}

type SystemBackupEncryptionReq struct {
	Format base.FileEncryptionFormat `json:"format,omitempty"`
	Secret string                    `json:"secret,omitzero"`
}

func (req *SystemBackupEncryptionReq) ToEntity() entity.SystemBackupEncryption {
	return entity.SystemBackupEncryption{
		Format: req.Format,
		Secret: entity.NewEncryptedField(req.Secret),
	}
}

type SystemBackupCloudStorageReq struct {
	ID             string `json:"id"`
	DestinationDir string `json:"destinationDir"`
}

func (req *SystemBackupCloudStorageReq) ToEntity() entity.SystemBackupCloudStorage {
	return entity.SystemBackupCloudStorage{
		ID:             req.ID,
		DestinationDir: req.DestinationDir,
	}
}

type SystemBackupDBConfigReq struct {
	BackupDeletedObjects bool `json:"backupDeletedObjects"`
}

func (req *SystemBackupDBConfigReq) ToEntity() entity.SystemBackupDBConfig {
	return entity.SystemBackupDBConfig{
		BackupDeletedObjects: req.BackupDeletedObjects,
	}
}

func (req *SystemBackupBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewUpdateSystemBackupReq() *UpdateSystemBackupReq {
	return &UpdateSystemBackupReq{}
}

func (req *UpdateSystemBackupReq) ModifyRequest() error {
	if !req.ScheduleFrom.IsZero() {
		req.ScheduleFrom = req.ScheduleFrom.Truncate(time.Minute)
	}
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *UpdateSystemBackupReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateSystemBackupResp struct {
	Meta *basedto.Meta `json:"meta"`
}

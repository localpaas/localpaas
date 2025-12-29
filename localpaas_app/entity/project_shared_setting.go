package entity

import (
	"time"
)

var (
	ProjectSharedSettingUpsertingConflictCols = []string{"project_id", "setting_id"}
	ProjectSharedSettingUpsertingUpdateCols   = []string{"data_view_allowed", "deleted_at"}
)

type ProjectSharedSetting struct {
	ProjectID       string `bun:",pk"`
	SettingID       string `bun:",pk"`
	DataViewAllowed bool

	CreatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

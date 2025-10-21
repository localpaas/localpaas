package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	SettingUpsertingConflictCols = []string{"id"}
	SettingUpsertingUpdateCols   = []string{"target_type", "target_id", "data",
		"updated_at", "updated_by", "deleted_at"}
)

type Setting struct {
	ID         string `bun:",pk"`
	TargetType base.SettingTargetType
	TargetID   string `bun:",nullzero"`
	Data       []byte

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	CreatedByUser *User `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (u *Setting) GetID() string {
	return u.ID
}

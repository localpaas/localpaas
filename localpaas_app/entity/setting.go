package entity

import (
	"encoding/json"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/pkg/reflectutil"
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
	Data       string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	CreatedByUser *User `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (s *Setting) GetID() string {
	return s.ID
}

func (s *Setting) SetData(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	s.Data = reflectutil.UnsafeBytesToStr(b)
	return nil
}

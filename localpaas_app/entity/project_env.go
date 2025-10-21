package entity

import (
	"time"
)

var (
	ProjectEnvUpsertingConflictCols = []string{"id"}
	ProjectEnvUpsertingUpdateCols   = []string{"name", "project_id", "display_order",
		"updated_at", "updated_by", "deleted_at"}
)

type ProjectEnv struct {
	ID           string `bun:",pk"`
	Name         string
	ProjectID    string
	DisplayOrder int

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	CreatedByUser *User `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (p *ProjectEnv) GetID() string {
	return p.ID
}

// GetName implements NamedEntity interface
func (p *ProjectEnv) GetName() string {
	return p.Name
}

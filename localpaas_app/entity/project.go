package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ProjectUpsertingConflictCols = []string{"id"}
	ProjectUpsertingUpdateCols   = []string{"name", "photo", "data", "status",
		"updated_at", "updated_by", "deleted_at"}
)

type Project struct {
	ID     string `bun:",pk"`
	Name   string
	Photo  string `bun:",nullzero"`
	Data   []byte `bun:",nullzero"`
	Status base.ProjectStatus

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	AllSettings   []*Setting    `bun:"rel:has-many,join:id=target_id"`
	Settings      []*Setting    `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:project"`
	Envs          []*ProjectEnv `bun:"rel:has-many,join:id=project_id"`
	Tags          []*ProjectTag `bun:"rel:has-many,join:id=project_id"`
	CreatedByUser *User         `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User         `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (p *Project) GetID() string {
	return p.ID
}

// GetName implements NamedEntity interface
func (p *Project) GetName() string {
	return p.Name
}

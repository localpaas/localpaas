package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	AppUpsertingConflictCols = []string{"id"}
	AppUpsertingUpdateCols   = []string{"name", "photo", "project_id", "project_env_id",
		"data", "status", "updated_at", "updated_by", "deleted_at"}
)

type App struct {
	ID           string `bun:",pk"`
	Name         string
	Photo        string `bun:",nullzero"`
	ProjectID    string
	ProjectEnvID string
	Data         []byte `bun:",nullzero"`
	Status       base.AppStatus

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Project       *Project    `bun:"rel:has-one,join:project_id=id"`
	ProjectEnv    *ProjectEnv `bun:"rel:has-one,join:project_env_id=id"`
	Settings      []*Setting  `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:app"`
	CreatedByUser *User       `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User       `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (p *App) GetID() string {
	return p.ID
}

// GetName implements NamedEntity interface
func (p *App) GetName() string {
	return p.Name
}

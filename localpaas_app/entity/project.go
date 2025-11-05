package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ProjectUpsertingConflictCols = []string{"id"}
	ProjectUpsertingUpdateCols   = []string{"name", "slug", "photo", "status", "note",
		"settings_id", "env_vars_id", "updated_at", "deleted_at"}
)

type Project struct {
	ID         string `bun:",pk"`
	Name       string
	Slug       string
	Photo      string `bun:",nullzero"`
	Status     base.ProjectStatus
	Note       string `bun:",nullzero"`
	SettingsID string `bun:",nullzero"`
	EnvVarsID  string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Settings *Setting         `bun:"rel:has-one,join:settings_id=id"`
	EnvVars  *Setting         `bun:"rel:has-one,join:env_vars_id=id"`
	Apps     []*App           `bun:"rel:has-many,join:id=project_id"`
	Tags     []*ProjectTag    `bun:"rel:has-many,join:id=project_id"`
	Accesses []*ACLPermission `bun:"rel:has-many,join:id=resource_id"`
}

// GetID implements IDEntity interface
func (p *Project) GetID() string {
	return p.ID
}

// GetName implements NamedEntity interface
func (p *Project) GetName() string {
	return p.Name
}

type ProjectSettings struct {
	Test string `json:"test"`
}

func (p *Project) ParseSettings() (*ProjectSettings, error) {
	if p.Settings != nil {
		res := &ProjectSettings{}
		return res, p.Settings.parseData(res)
	}
	return nil, nil
}

type ProjectEnvVars struct {
	Data [][]string `json:"data"`
}

func (p *Project) ParseEnvVars() (*ProjectEnvVars, error) {
	if p.EnvVars != nil {
		res := &ProjectEnvVars{}
		return res, p.EnvVars.parseData(res)
	}
	return nil, nil
}

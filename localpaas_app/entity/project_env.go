package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ProjectEnvUpsertingConflictCols = []string{"id"}
	ProjectEnvUpsertingUpdateCols   = []string{"name", "project_id", "display_order", "data", "status",
		"updated_at", "updated_by", "deleted_at"}
)

type ProjectEnv struct {
	ID           string `bun:",pk"`
	Name         string
	ProjectID    string
	DisplayOrder int
	Status       base.ProjectStatus

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	MainSettings    []*Setting `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:project-env"`
	EnvVarsSettings []*Setting `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:env-var"`
	Apps            []*App     `bun:"rel:has-many,join:id=project_env_id"`
	CreatedByUser   *User      `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser   *User      `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (p *ProjectEnv) GetID() string {
	return p.ID
}

// GetName implements NamedEntity interface
func (p *ProjectEnv) GetName() string {
	return p.Name
}

type ProjectEnvSettings struct {
	Test string `json:"test"`
}

func (p *ProjectEnv) GetMainSettings() (*ProjectEnvSettings, error) {
	if len(p.MainSettings) > 0 {
		res := &ProjectEnvSettings{}
		return res, p.MainSettings[0].parseData(res)
	}
	return nil, nil
}

type ProjectEnvEnvVars struct {
	Data [][]string `json:"data"`
}

func (p *ProjectEnv) GetEnvVars() (*ProjectEnvEnvVars, error) {
	if len(p.EnvVarsSettings) > 0 {
		res := &ProjectEnvEnvVars{}
		return res, p.EnvVarsSettings[0].parseData(res)
	}
	return nil, nil
}

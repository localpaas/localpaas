package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ProjectUpsertingConflictCols = []string{"id"}
	ProjectUpsertingUpdateCols   = []string{"name", "photo", "status", "note",
		"updated_at", "updated_by", "deleted_at"}
)

type Project struct {
	ID     string `bun:",pk"`
	Name   string
	Photo  string `bun:",nullzero"`
	Status base.ProjectStatus
	Note   string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	EnvVarsSettings []*Setting    `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:env-var"`
	MainSettings    []*Setting    `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:project"`
	Apps            []*App        `bun:"rel:has-many,join:id=project_id"`
	Tags            []*ProjectTag `bun:"rel:has-many,join:id=project_id"`
	CreatedByUser   *User         `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser   *User         `bun:"rel:has-one,join:updated_by=id"`
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

func (p *Project) GetMainSettings() (*ProjectSettings, error) {
	if len(p.MainSettings) > 0 {
		res := &ProjectSettings{}
		return res, p.MainSettings[0].parseData(res)
	}
	return nil, nil
}

type ProjectEnvVars struct {
	Data [][]string `json:"data"`
}

func (p *Project) GetEnvVars() (*ProjectEnvVars, error) {
	if len(p.EnvVarsSettings) > 0 {
		res := &ProjectEnvVars{}
		return res, p.EnvVarsSettings[0].parseData(res)
	}
	return nil, nil
}

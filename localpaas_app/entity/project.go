package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ProjectUpsertingConflictCols = []string{"id"}
	ProjectUpsertingUpdateCols   = []string{"name", "slug", "photo", "status", "note",
		"updated_at", "deleted_at"}
)

type Project struct {
	ID     string `bun:",pk"`
	Name   string
	Slug   string
	Photo  string `bun:",nullzero"`
	Status base.ProjectStatus
	Note   string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Settings []*Setting       `bun:"rel:has-many,join:id=object_id"`
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

func (p *Project) GetSettingsByType(typ base.SettingType) (resp []*Setting) {
	for _, setting := range p.Settings {
		if setting.Type == typ {
			resp = append(resp, setting)
		}
	}
	return resp
}

func (p *Project) GetSettingByType(typ base.SettingType) *Setting {
	for _, setting := range p.Settings {
		if setting.Type == typ {
			return setting
		}
	}
	return nil
}

func (p *Project) GetDefaultNetworkName() string {
	return p.Slug + "_net"
}

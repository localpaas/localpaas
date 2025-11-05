package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	AppUpsertingConflictCols = []string{"id"}
	AppUpsertingUpdateCols   = []string{"name", "slug", "photo", "project_id", "parent_id", "status", "note",
		"settings_id", "env_vars_id", "updated_at", "deleted_at"}
)

type App struct {
	ID         string `bun:",pk"`
	Name       string
	Slug       string
	Photo      string `bun:",nullzero"`
	ProjectID  string
	ParentID   string `bun:",nullzero"`
	Status     base.AppStatus
	Note       string `bun:",nullzero"`
	SettingsID string `bun:",nullzero"`
	EnvVarsID  string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Project  *Project         `bun:"rel:has-one,join:project_id=id"`
	Parent   *App             `bun:"rel:has-one,join:parent_id=id"`
	Settings *Setting         `bun:"rel:has-one,join:settings_id=id"`
	EnvVars  *Setting         `bun:"rel:has-one,join:env_vars_id=id"`
	Tags     []*AppTag        `bun:"rel:has-many,join:id=app_id"`
	Accesses []*ACLPermission `bun:"rel:has-many,join:id=resource_id"`
}

// GetID implements IDEntity interface
func (app *App) GetID() string {
	return app.ID
}

// GetName implements NamedEntity interface
func (app *App) GetName() string {
	return app.Name
}

type AppSettings struct {
	Test string `json:"test"`
}

func (app *App) ParseSettings() (*AppSettings, error) {
	if app.Settings != nil {
		res := &AppSettings{}
		return res, app.Settings.parseData(res)
	}
	return nil, nil
}

type AppEnvVars struct {
	Data [][]string `json:"data"`
}

func (app *App) ParseEnvVars() (*AppEnvVars, error) {
	if app.EnvVars != nil {
		res := &AppEnvVars{}
		return res, app.EnvVars.parseData(res)
	}
	return nil, nil
}

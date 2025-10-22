package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	AppUpsertingConflictCols = []string{"id"}
	AppUpsertingUpdateCols   = []string{"name", "photo", "project_id", "project_env_id",
		"status", "updated_at", "updated_by", "deleted_at"}
)

type App struct {
	ID           string `bun:",pk"`
	Name         string
	Photo        string `bun:",nullzero"`
	ProjectID    string
	ProjectEnvID string
	Status       base.AppStatus

	CreatedAt time.Time `bun:",default:current_timestamp"`
	CreatedBy string
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedBy string
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Project         *Project    `bun:"rel:has-one,join:project_id=id"`
	ProjectEnv      *ProjectEnv `bun:"rel:has-one,join:project_env_id=id"`
	MainSettings    []*Setting  `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:app"`
	EnvVarsSettings []*Setting  `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:env-var"`
	Tags            []*AppTag   `bun:"rel:has-many,join:id=app_id"`
	CreatedByUser   *User       `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser   *User       `bun:"rel:has-one,join:updated_by=id"`
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

func (app *App) GetMainSettings() (*AppSettings, error) {
	if len(app.MainSettings) > 0 {
		res := &AppSettings{}
		return res, app.MainSettings[0].parseData(res)
	}
	return nil, nil
}

type AppEnvVars struct {
	Data [][]string `json:"data"`
}

func (app *App) GetEnvVars() (*AppEnvVars, error) {
	if len(app.EnvVarsSettings) > 0 {
		res := &AppEnvVars{}
		return res, app.EnvVarsSettings[0].parseData(res)
	}
	return nil, nil
}

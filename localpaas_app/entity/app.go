package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	AppUpsertingConflictCols = []string{"id"}
	AppUpsertingUpdateCols   = []string{"name", "key", "project_id", "parent_id", "service_id",
		"status", "token", "webhook_secret", "note", "update_ver", "updated_at", "deleted_at"}
	AppDefaultExcludeColumns = []string{"note"}
)

type App struct {
	ID            string `bun:",pk"`
	Name          string
	Key           string
	ProjectID     string
	ParentID      string `bun:",nullzero"`
	ServiceID     string `bun:",nullzero"`
	Status        base.AppStatus
	Token         string
	WebhookSecret string `bun:",nullzero"`
	Note          string `bun:",nullzero"`
	UpdateVer     int

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Project   *Project   `bun:"rel:has-one,join:project_id=id"`
	ParentApp *App       `bun:"rel:has-one,join:parent_id=id"`
	Settings  []*Setting `bun:"rel:has-many,join:id=object_id"`
	Tags      []*AppTag  `bun:"rel:has-many,join:id=app_id"`
}

// GetID implements IDEntity interface
func (app *App) GetID() string {
	return app.ID
}

// GetName implements NamedEntity interface
func (app *App) GetName() string {
	return app.Name
}

func (app *App) GetSettingsByType(typ base.SettingType) (resp []*Setting) {
	for _, setting := range app.Settings {
		if setting.Type == typ {
			resp = append(resp, setting)
		}
	}
	return resp
}

func (app *App) GetSettingByType(typ base.SettingType) *Setting {
	for _, setting := range app.Settings {
		if setting.Type == typ {
			return setting
		}
	}
	return nil
}

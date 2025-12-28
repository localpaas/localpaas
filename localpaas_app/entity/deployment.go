package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentDeploymentVersion = 1
)

var (
	DeploymentUpsertingConflictCols = []string{"id"}
	DeploymentUpsertingUpdateCols   = []string{"app_id", "deployment_settings", "status",
		"version", "update_ver", "started_at", "ended_at", "updated_at", "deleted_at"}
)

type Deployment struct {
	ID                 string `bun:",pk"`
	AppID              string
	DeploymentSettings *AppDeploymentSettings
	Status             base.DeploymentStatus
	Version            int
	UpdateVer          int

	StartedAt time.Time `bun:",nullzero"`
	EndedAt   time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	App *App `bun:"rel:belongs-to,join:app_id=id"`
}

// GetID implements IDEntity interface
func (d *Deployment) GetID() string {
	return d.ID
}

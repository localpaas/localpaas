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
	DeploymentUpsertingUpdateCols   = []string{"app_id", "status", "settings", "output",
		"version", "update_ver", "started_at", "ended_at", "updated_at", "deleted_at"}
)

type Deployment struct {
	ID        string `bun:",pk"`
	AppID     string
	Status    base.DeploymentStatus
	Settings  *AppDeploymentSettings
	Output    *AppDeploymentOutput
	Version   int
	UpdateVer int

	StartedAt time.Time `bun:",nullzero"`
	EndedAt   time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	App *App `bun:"rel:belongs-to,join:app_id=id"`
}

type AppDeploymentOutput struct {
	CommitHash    string   `json:"commitHash,omitempty"`
	CommitMessage string   `json:"commitMessage,omitempty"`
	ImageTags     []string `json:"imageTags,omitempty"`
}

// GetID implements IDEntity interface
func (d *Deployment) GetID() string {
	return d.ID
}

func (d *Deployment) CanCancel() bool {
	if d.Status == base.DeploymentStatusDone ||
		d.Status == base.DeploymentStatusCanceled ||
		d.Status == base.DeploymentStatusFailed {
		return false
	}
	return true
}

package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	TaskUpsertingConflictCols = []string{"id"}
	TaskUpsertingUpdateCols   = []string{"job_id", "type", "status", "data", "error", "version",
		"run_at", "done_at", "updated_at", "deleted_at"}
)

type Task struct {
	ID      string `bun:",pk"`
	JobID   string `bun:",nullzero"`
	Type    base.TaskType
	Status  base.TaskStatus
	Data    string `bun:",nullzero"`
	Error   string `bun:",nullzero"`
	Version int

	RunAt  time.Time `bun:",nullzero"`
	DoneAt time.Time `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Job *Setting `bun:"rel:belongs-to,join:job_id=id"`
}

// GetID implements IDEntity interface
func (t *Task) GetID() string {
	return t.ID
}

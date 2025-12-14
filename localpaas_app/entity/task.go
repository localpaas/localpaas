package entity

import (
	"encoding/json"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

const (
	CurrentTaskVersion = 1
)

var (
	TaskUpsertingConflictCols = []string{"id"}
	TaskUpsertingUpdateCols   = []string{"job_id", "type", "status", "args",
		"priority", "max_retry", "retry", "retry_delay_secs", "runs", "version",
		"run_at", "started_at", "ended_at", "updated_at", "deleted_at"}
)

type Task struct {
	ID             string `bun:",pk"`
	JobID          string
	Type           base.TaskType
	Status         base.TaskStatus
	Priority       base.TaskPriority
	MaxRetry       int
	Retry          int
	RetryDelaySecs int
	Args           string
	Runs           string
	Version        int

	RunAt     time.Time
	StartedAt time.Time `bun:",nullzero"`
	EndedAt   time.Time `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	Job          *Setting      `bun:"rel:belongs-to,join:job_id=id"`
	UpdatingTask *UpdatingTask `bun:"rel:has-one,join:id=id"`
}

// GetID implements IDEntity interface
func (t *Task) GetID() string {
	return t.ID
}

func (t *Task) GetRuns() ([]*TaskRun, error) {
	if t.Runs == "" {
		return nil, nil
	}
	runs := []*TaskRun{}
	err := json.Unmarshal(reflectutil.UnsafeStrToBytes(t.Runs), &runs)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return runs, nil
}

func (t *Task) AddRun(run *TaskRun) error {
	runs, err := t.GetRuns()
	if err != nil {
		return apperrors.Wrap(err)
	}

	runs = append(runs, run)
	runBytes, err := json.Marshal(runs)
	if err != nil {
		return apperrors.Wrap(err)
	}
	t.Runs = reflectutil.UnsafeBytesToStr(runBytes)
	return nil
}

type TaskRun struct {
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
	Error     string    `json:"error,omitempty"`
}

var (
	UpdatingTaskUpsertingConflictCols = []string{"id"}
	UpdatingTaskUpsertingUpdateCols   = []string{"started_at"}
)

type UpdatingTask struct {
	ID        string `bun:",pk"`
	StartedAt time.Time
}

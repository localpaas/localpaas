package gocronqueue

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	taskQueueSchedKey         = "task:queue:sched"
	taskQueueSchedReadTimeout = 10 * time.Minute
)

type SchedMessage struct {
	SchedTasks     []*entity.Task `json:"schedTasks,omitempty"`
	UnschedTaskIDs []string       `json:"unschedTaskIDs,omitempty"`
}

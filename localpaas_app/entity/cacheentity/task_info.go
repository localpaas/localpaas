package cacheentity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type TaskInfo struct {
	ID        string          `json:"id"`
	Status    base.TaskStatus `json:"status"`
	StartedAt time.Time       `json:"startedAt"`
}

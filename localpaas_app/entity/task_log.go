package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
)

type TaskLog struct {
	ID       int64 `bun:",pk,autoincrement"`
	TaskID   string
	TargetID string         `bun:",nullzero"`
	Type     applog.LogType `bun:",nullzero"`
	Data     string
	Ts       time.Time `bun:",nullzero"`
}

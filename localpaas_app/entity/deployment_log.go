package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/pkg/realtimelog"
)

type DeploymentLog struct {
	ID           int64 `bun:",pk,autoincrement"`
	DeploymentID string
	Type         realtimelog.LogType `bun:",nullzero"`
	Data         string
	Ts           time.Time `bun:",nullzero"`
}

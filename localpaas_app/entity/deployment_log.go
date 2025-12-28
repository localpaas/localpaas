package entity

import (
	"time"
)

type DeploymentLog struct {
	ID           int64 `bun:",pk,autoincrement"`
	DeploymentID string
	Step         string `bun:",nullzero"`
	Content      string
	CreatedAt    time.Time `bun:",default:current_timestamp"`
}

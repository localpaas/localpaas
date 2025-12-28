package cacheentity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type DeploymentInfo struct {
	ID        string                `json:"id"`
	AppID     string                `json:"appId"`
	Status    base.DeploymentStatus `json:"status"`
	Cancel    bool                  `json:"cancel,omitempty"`
	StartedAt time.Time             `json:"startedAt"`
}

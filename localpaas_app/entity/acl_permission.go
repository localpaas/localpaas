package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ACLPermissionUpsertingConflictCols = []string{"user_id", "resource_type", "resource_id"}
	ACLPermissionUpsertingUpdateCols   = []string{"action_read", "action_write", "action_delete",
		"updated_at", "updated_by"}
)

type ACLPermission struct {
	UserID       string            `bun:",pk"`
	ResourceType base.ResourceType `bun:",pk"`
	ResourceID   string            `bun:",pk"`
	Actions      AccessActions     `bun:"embed:action_"`
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string

	User          *User `bun:"rel:has-one,join:user_id=id"`
	CreatedByUser *User `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User `bun:"rel:has-one,join:updated_by=id"`
}

type AccessActions struct {
	Read   base.AccessType `json:"read"`
	Write  base.AccessType `json:"write"`
	Delete base.AccessType `json:"delete"`
}

func (a *AccessActions) Equal(other AccessActions) bool {
	return a.Read == other.Read && a.Write == other.Write && a.Delete == other.Delete
}

func (a *AccessActions) IsFullAccess() bool {
	return a.Read == base.AccessTypeYes && a.Write == base.AccessTypeYes && a.Delete == base.AccessTypeYes
}

func (a *AccessActions) IsNoAccess() bool {
	return a.Read == base.AccessTypeNo && a.Write == base.AccessTypeNo && a.Delete == base.AccessTypeNo
}

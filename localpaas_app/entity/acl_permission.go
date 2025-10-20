package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ACLPermissionUpsertingConflictCols = []string{"id"}
	ACLPermissionUpsertingUpdateCols   = []string{"action_read", "action_write", "action_delete",
		"updated_at"}
)

type ACLPermission struct {
	ID           string `bun:",pk"`
	UserID       string
	ResourceType base.ResourceType
	ResourceID   string
	Actions      AccessActions `bun:"embed:action_"`
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time

	CreatedByUser *User `bun:"rel:has-one,join:created_by=id"`
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

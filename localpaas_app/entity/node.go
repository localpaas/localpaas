package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	NodeUpsertingConflictCols = []string{"id"}
	NodeUpsertingUpdateCols   = []string{"is_leader", "is_manager", "host_name", "ip", "status", "infra_status",
		"info", "note", "last_synced_at", "updated_at", "updated_by", "deleted_at"}
)

type Node struct {
	ID          string `bun:",pk"`
	IsLeader    bool
	IsManager   bool
	HostName    string
	IP          string
	Status      base.NodeStatus
	InfraStatus string
	Info        string `bun:",nullzero"`
	Note        string `bun:",nullzero"`

	LastSyncedAt time.Time
	CreatedAt    time.Time `bun:",default:current_timestamp"`
	CreatedBy    string
	UpdatedAt    time.Time `bun:",default:current_timestamp"`
	UpdatedBy    string
	DeletedAt    time.Time `bun:",soft_delete,nullzero"`

	MainSettings  []*Setting `bun:"rel:has-many,join:id=target_id,join:type=target_type,polymorphic:node"`
	CreatedByUser *User      `bun:"rel:has-one,join:created_by=id"`
	UpdatedByUser *User      `bun:"rel:has-one,join:updated_by=id"`
}

// GetID implements IDEntity interface
func (n *Node) GetID() string {
	return n.ID
}

// GetName implements NamedEntity interface
func (n *Node) GetName() string {
	return n.HostName
}

type NodeSettings struct {
	Test string `json:"test"`
}

func (n *Node) GetMainSettings() (*NodeSettings, error) {
	if len(n.MainSettings) > 0 {
		res := &NodeSettings{}
		return res, n.MainSettings[0].parseData(res)
	}
	return nil, nil
}

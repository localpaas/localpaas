package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	NodeUpsertingConflictCols = []string{"id"}
	NodeUpsertingUpdateCols   = []string{"is_leader", "is_manager", "host_name", "ip", "status", "infra_status",
		"info", "note", "settings_id", "last_synced_at", "updated_at", "deleted_at"}
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
	SettingsID  string `bun:",nullzero"`

	LastSyncedAt time.Time
	CreatedAt    time.Time `bun:",default:current_timestamp"`
	UpdatedAt    time.Time `bun:",default:current_timestamp"`
	DeletedAt    time.Time `bun:",soft_delete,nullzero"`

	Settings *Setting `bun:"rel:has-one,join:settings_id=id"`
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

func (n *Node) ParseSettings() (*NodeSettings, error) {
	if n.Settings != nil {
		res := &NodeSettings{}
		return res, n.Settings.parseData(res)
	}
	return nil, nil
}

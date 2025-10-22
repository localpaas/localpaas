package base

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"
	NodeStatusInactive NodeStatus = "inactive"
	NodeStatusDeleting NodeStatus = "deleting"
)

var (
	AllNodeStatuses = []NodeStatus{NodeStatusActive, NodeStatusInactive, NodeStatusDeleting}
)

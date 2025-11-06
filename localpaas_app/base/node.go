package base

import "github.com/docker/docker/api/types/swarm"

type NodeStatus string

const (
	NodeStatusUnknown      = NodeStatus(swarm.NodeStateUnknown)
	NodeStatusDown         = NodeStatus(swarm.NodeStateDown)
	NodeStatusReady        = NodeStatus(swarm.NodeStateReady)
	NodeStatusDisconnected = NodeStatus(swarm.NodeStateDisconnected)
)

var (
	AllNodeStatuses = []NodeStatus{NodeStatusUnknown, NodeStatusDown, NodeStatusReady, NodeStatusDisconnected}
)

type NodeRole string

const (
	NodeRoleManager = NodeRole(swarm.NodeRoleManager)
	NodeRoleWorker  = NodeRole(swarm.NodeRoleWorker)
)

var (
	AllNodeRoles = []NodeRole{NodeRoleManager, NodeRoleWorker}
)

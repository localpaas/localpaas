package base

type ResourceType string

const (
	ResourceTypeUser       ResourceType = "user"
	ResourceTypeCluster    ResourceType = "cluster"
	ResourceTypeNode       ResourceType = "node"
	ResourceTypeNetwork    ResourceType = "network"
	ResourceTypeDeployment ResourceType = "deployment"
	ResourceTypeProject    ResourceType = "project"
	ResourceTypeApp        ResourceType = "app"
)

type ActionType string

const (
	ActionTypeRead    ActionType = "read"
	ActionTypeWrite   ActionType = "write"
	ActionTypeDelete  ActionType = "delete"
	ActionTypeExecute ActionType = "execute"
)

var (
	AllActionTypes = []ActionType{ActionTypeRead, ActionTypeWrite, ActionTypeDelete, ActionTypeExecute}
)

type AccessType string

const (
	AccessTypeYes = AccessType("yes")
	AccessTypeNo  = AccessType("no")
)

var (
	AllAccessTypes = []AccessType{AccessTypeYes, AccessTypeNo}
)

type PermissionResource struct {
	UserID       string
	ResourceType ResourceType
	ResourceID   string
}

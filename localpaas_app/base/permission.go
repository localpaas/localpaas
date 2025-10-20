package base

type ResourceType string

const (
	ResourceTypeUser    ResourceType = "user"
	ResourceTypeProject ResourceType = "project"
)

type ActionType string

const (
	ActionTypeRead   ActionType = "read"
	ActionTypeWrite  ActionType = "write"
	ActionTypeDelete ActionType = "delete"
)

var (
	AllActionTypes = []ActionType{ActionTypeRead, ActionTypeWrite, ActionTypeDelete}
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

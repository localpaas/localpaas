package base

type SubjectType string

const (
	SubjectTypeUser       SubjectType = "user"
	SubjectTypeCluster    SubjectType = "cluster"
	SubjectTypeNode       SubjectType = "node"
	SubjectTypeNetwork    SubjectType = "network"
	SubjectTypeDeployment SubjectType = "deployment"
	SubjectTypeProject    SubjectType = "project"
	SubjectTypeApp        SubjectType = "app"
)

type ResourceType string

const (
	ResourceTypeModule     ResourceType = "module"
	ResourceTypeUser       ResourceType = "user"
	ResourceTypeCluster    ResourceType = "cluster"
	ResourceTypeNode       ResourceType = "node"
	ResourceTypeNetwork    ResourceType = "network"
	ResourceTypeDeployment ResourceType = "deployment"
	ResourceTypeProject    ResourceType = "project"
	ResourceTypeApp        ResourceType = "app"
	ResourceTypeS3Storage  ResourceType = "s3-storage"
	ResourceTypeSSHKey     ResourceType = "ssh-key"
	ResourceTypeAPIKey     ResourceType = "api-key"
)

type ResourceModule string

const (
	ResourceModuleUserAndRole ResourceModule = "user-role"
	ResourceModuleSetting     ResourceModule = "setting"
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

type PermissionResource struct {
	SubjectType  SubjectType
	SubjectID    string
	ResourceType ResourceType
	ResourceID   string
}

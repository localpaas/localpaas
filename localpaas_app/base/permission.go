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
	ResourceTypeModule       ResourceType = "module"
	ResourceTypeUser         ResourceType = "user"
	ResourceTypeCluster      ResourceType = "cluster"
	ResourceTypeNode         ResourceType = "node"
	ResourceTypeNetwork      ResourceType = "network"
	ResourceTypeVolume       ResourceType = "volume"
	ResourceTypeImage        ResourceType = "image"
	ResourceTypeDeployment   ResourceType = "deployment"
	ResourceTypeProject      ResourceType = "project"
	ResourceTypeApp          ResourceType = "app"
	ResourceTypeS3Storage    ResourceType = "s3-storage"
	ResourceTypeOAuth        ResourceType = "oauth"
	ResourceTypeSSHKey       ResourceType = "ssh-key"
	ResourceTypeAPIKey       ResourceType = "api-key"
	ResourceTypeSecret       ResourceType = "secret"
	ResourceTypeSlack        ResourceType = "slack"
	ResourceTypeDiscord      ResourceType = "discord"
	ResourceTypeRegistryAuth ResourceType = "registry-auth"
	ResourceTypeBasicAuth    ResourceType = "basic-auth"
	ResourceTypeSsl          ResourceType = "ssl"
)

type ResourceModule string

const (
	ResourceModuleSettings ResourceModule = "mod::settings"
	ResourceModuleProvider ResourceModule = "mod::provider"
	ResourceModuleCluster  ResourceModule = "mod::cluster"
	ResourceModuleUser     ResourceModule = "mod::user"
	ResourceModuleProject  ResourceModule = "mod::project"
)

var (
	AllResourceModules = []ResourceModule{ResourceModuleSettings, ResourceModuleProvider, ResourceModuleUser,
		ResourceModuleCluster, ResourceModuleProject}
)

type ActionType string

const (
	ActionTypeRead   ActionType = "read"
	ActionTypeWrite  ActionType = "write"
	ActionTypeDelete ActionType = "delete"
)

var (
	AllActionTypes = []ActionType{ActionTypeRead, ActionTypeWrite, ActionTypeDelete}

	//nolint
	mapActionValues = map[ActionType]int{
		ActionTypeRead:   1,
		ActionTypeWrite:  2,
		ActionTypeDelete: 3,
	}
)

func ActionTypeCmp(a1, s2 ActionType) int {
	return mapActionValues[a1] - mapActionValues[s2]
}

type PermissionResource struct {
	SubjectType  SubjectType
	SubjectID    string
	ResourceType ResourceType
	ResourceID   string
}

package permission

import "github.com/localpaas/localpaas/localpaas_app/base"

type AccessCheck struct {
	SubjectType        base.SubjectType
	SubjectID          string
	ResourceModule     base.ResourceModule
	ResourceType       base.ResourceType
	ResourceID         string
	ParentResourceType base.ResourceType
	ParentResourceID   string
	Action             base.ActionType
}

package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

var (
	ACLPermissionUpsertingConflictCols = []string{"subject_id", "resource_id"}
	ACLPermissionUpsertingUpdateCols   = []string{"subject_type", "resource_type",
		"action_read", "action_write", "action_delete", "updated_at", "deleted_at"}
)

type ACLPermission struct {
	SubjectType  base.SubjectType   `json:"subjectType"`
	SubjectID    string             `bun:",pk" json:"subjectId"`
	ResourceType base.ResourceType  `json:"resourceType"`
	ResourceID   string             `bun:",pk" json:"resourceId"`
	Actions      base.AccessActions `bun:"embed:action_" json:"actions"`

	CreatedAt time.Time `bun:",default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:",default:current_timestamp" json:"updatedAt"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deletedAt,omitzero"`

	SubjectUser    *User    `bun:"rel:has-one,join:subject_id=id" json:"subjectUser,omitempty"`
	SubjectProject *Project `bun:"rel:has-one,join:subject_id=id" json:"subjectProject,omitempty"`
	SubjectApp     *App     `bun:"rel:has-one,join:subject_id=id" json:"subjectApp,omitempty"`

	ResourceProject *Project `bun:"rel:has-one,join:resource_id=id" json:"resourceProject,omitempty"`
	ResourceApp     *App     `bun:"rel:has-one,join:resource_id=id" json:"resourceApp,omitempty"`
}

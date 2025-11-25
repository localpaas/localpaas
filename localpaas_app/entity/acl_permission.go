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
	SubjectType  base.SubjectType
	SubjectID    string `bun:",pk"`
	ResourceType base.ResourceType
	ResourceID   string             `bun:",pk"`
	Actions      base.AccessActions `bun:"embed:action_"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	SubjectUser    *User    `bun:"rel:has-one,join:subject_id=id"`
	SubjectProject *Project `bun:"rel:has-one,join:subject_id=id"`
	SubjectApp     *App     `bun:"rel:has-one,join:subject_id=id"`

	ResourceProject *Project `bun:"rel:has-one,join:resource_id=id"`
	ResourceApp     *App     `bun:"rel:has-one,join:resource_id=id"`
}

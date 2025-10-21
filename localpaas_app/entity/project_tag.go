package entity

import (
	"time"
)

var (
	ProjectTagUpsertingConflictCols = []string{"project_id", "tag"}
	ProjectTagUpsertingUpdateCols   = []string{"display_order", "deleted_at"}
)

type ProjectTag struct {
	ProjectID    string `bun:",pk"`
	Tag          string `bun:",pk"`
	DisplayOrder int

	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

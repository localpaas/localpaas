package entity

import (
	"time"
)

var (
	AppTagUpsertingConflictCols = []string{"app_id", "tag"}
	AppTagUpsertingUpdateCols   = []string{"display_order", "deleted_at"}
)

type AppTag struct {
	AppID        string `bun:",pk"`
	Tag          string `bun:",pk"`
	DisplayOrder int

	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

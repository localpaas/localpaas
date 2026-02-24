package entity

var (
	LockUpsertingConflictCols = []string{"id"}
	LockUpsertingUpdateCols   = []string{}
)

type Lock struct {
	ID string `bun:",pk"`
}

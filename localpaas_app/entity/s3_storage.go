package entity

import (
	"time"
)

var (
	S3StorageUpsertingConflictCols = []string{"id"}
	S3StorageUpsertingUpdateCols   = []string{"name", "access_key_id", "secret_access_key", "salt",
		"region", "bucket", "updated_at", "deleted_at"}
)

type S3Storage struct {
	ID              string `bun:",pk"`
	Name            string
	AccessKeyID     string
	SecretAccessKey []byte
	Salt            []byte
	Region          string `bun:",nullzero"`
	Bucket          string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	ObjectAccesses []*ACLPermission `bun:"rel:has-many,join:id=resource_id"`
}

// GetID implements IDEntity interface
func (p *S3Storage) GetID() string {
	return p.ID
}

// GetName implements NamedEntity interface
func (p *S3Storage) GetName() string {
	return p.Name
}

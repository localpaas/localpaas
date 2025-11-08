package entity

import (
	"encoding/json"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/pkg/reflectutil"
)

var (
	SettingUpsertingConflictCols = []string{"id"}
	SettingUpsertingUpdateCols   = []string{"object_id", "type", "status", "name", "data",
		"updated_at", "expire_at", "deleted_at"}
)

type Setting struct {
	ID       string `bun:",pk"`
	ObjectID string `bun:",nullzero"`
	Type     base.SettingType
	Status   base.SettingStatus
	Name     string `bun:",nullzero"`
	Data     string `bun:",nullzero"`

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	ExpireAt  time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	ObjectAccesses []*ACLPermission `bun:"rel:has-many,join:id=resource_id"`
	ObjectUser     *User            `bun:"rel:belongs-to,join:object_id=id"`
	ObjectProject  *Project         `bun:"rel:belongs-to,join:object_id=id"`
	ObjectApp      *App             `bun:"rel:belongs-to,join:object_id=id"`
}

// GetID implements IDEntity interface
func (s *Setting) GetID() string {
	return s.ID
}

// GetName implements NamedEntity interface
func (s *Setting) GetName() string {
	return s.Name
}

func (s *Setting) parseData(structPtr any) error {
	if s == nil || s.Data == "" {
		return nil
	}
	err := json.Unmarshal(reflectutil.UnsafeStrToBytes(s.Data), structPtr)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *Setting) SetData(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	s.Data = reflectutil.UnsafeBytesToStr(b)
	return nil
}

func (s *Setting) MustSetData(data any) {
	gofn.Must1(s.SetData(data))
}

package entity

import (
	"encoding/json"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

var (
	SettingUpsertingConflictCols = []string{"id"}
	SettingUpsertingUpdateCols   = []string{"object_id", "type", "kind", "status", "name", "data",
		"avail_in_projects", "version", "update_ver", "updated_at", "expire_at", "deleted_at"}
)

type Setting struct {
	ID              string `bun:",pk"`
	ObjectID        string `bun:",nullzero"`
	Type            base.SettingType
	Kind            string `bun:",nullzero"`
	Status          base.SettingStatus
	Name            string `bun:",nullzero"`
	Data            string `bun:",nullzero"`
	AvailInProjects bool
	Version         int
	UpdateVer       int

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	ExpireAt  time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	ObjectAccesses []*ACLPermission `bun:"rel:has-many,join:id=resource_id"`
	ObjectUser     *User            `bun:"rel:belongs-to,join:object_id=id"`
	ObjectProject  *Project         `bun:"rel:belongs-to,join:object_id=id"`
	ObjectApp      *App             `bun:"rel:belongs-to,join:object_id=id"`
	Tasks          []*Task          `bun:"rel:has-many,join:id=job_id"`

	// NOTE: temporary fields
	parsedData      any
	CurrentObjectID string `bun:"-"`
}

// GetID implements IDEntity interface
func (s *Setting) GetID() string {
	return s.ID
}

// GetName implements NamedEntity interface
func (s *Setting) GetName() string {
	return s.Name
}

// IsActive returns true if setting has status `active` and does not expire
func (s *Setting) IsActive() bool {
	return s.Status == base.SettingStatusActive && !s.IsExpired()
}

func (s *Setting) IsExpired() bool {
	return !s.ExpireAt.IsZero() && s.ExpireAt.Before(time.Now())
}

func (s *Setting) IsStatusDirty() bool {
	return s.Status == base.SettingStatusActive && s.IsExpired()
}

func (s *Setting) parseData(structPtr any) error {
	if s == nil || s.Data == "" {
		return nil
	}
	err := json.Unmarshal(reflectutil.UnsafeStrToBytes(s.Data), structPtr)
	if err != nil {
		return apperrors.Wrap(err)
	}
	s.parsedData = structPtr
	return nil
}

func (s *Setting) SetData(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	s.Data = reflectutil.UnsafeBytesToStr(b)
	s.parsedData = data
	return nil
}

func (s *Setting) MustSetData(data any) {
	gofn.Must1(s.SetData(data))
}

func parseSettingAs[T any](s *Setting, typ base.SettingType, newFn func() T) (res T, error error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(T)
		if !ok {
			return res, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	if s.Data != "" && s.Type == typ {
		res = newFn()
		if err := s.parseData(res); err != nil {
			return res, apperrors.Wrap(err)
		}
	}
	return res, nil
}

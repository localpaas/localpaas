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
		"avail_in_projects", "is_default", "ref_ids", "version", "update_ver",
		"updated_at", "expire_at", "deleted_at"}
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
	Default         bool     `bun:"is_default"`
	RefIDs          []string `bun:"ref_ids,array,nullzero"`
	Version         int
	UpdateVer       int

	CreatedAt time.Time `bun:",default:current_timestamp"`
	UpdatedAt time.Time `bun:",default:current_timestamp"`
	ExpireAt  time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`

	ObjectUser *User   `bun:"rel:belongs-to,join:object_id=id"`
	Tasks      []*Task `bun:"rel:has-many,join:id=job_id"`

	// NOTE: temporary fields
	parsedData      SettingData
	CurrentObjectID string     `bun:"-"`
	RefSettings     []*Setting `bun:"-"`
}

type SettingData interface {
	GetType() base.SettingType
	GetRefSettingIDs() []string
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

func (s *Setting) parseData(structPtr SettingData) error {
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

func (s *Setting) SetData(data SettingData) error {
	if data.GetType() != s.Type {
		return apperrors.NewTypeInvalid()
	}
	b, err := json.Marshal(data)
	if err != nil {
		return apperrors.Wrap(err)
	}
	s.Data = reflectutil.UnsafeBytesToStr(b)
	s.parsedData = data
	s.RefIDs = s.parsedData.GetRefSettingIDs()
	return nil
}

func (s *Setting) MustSetData(data SettingData) {
	gofn.Must1(s.SetData(data))
}

func parseSettingAs[T SettingData](s *Setting, newFn func() T) (res T, error error) {
	if s.parsedData != nil {
		res, ok := s.parsedData.(T)
		if !ok {
			return res, apperrors.NewTypeInvalid()
		}
		return res, nil
	}
	if s.Data != "" {
		res = newFn()
		if res.GetType() != s.Type {
			return res, apperrors.NewTypeInvalid()
		}
		if err := s.parseData(res); err != nil {
			return res, apperrors.Wrap(err)
		}
	}
	return res, nil
}

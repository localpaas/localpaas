package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

var (
	UserUpsertingConflictCols = []string{"id"}
	UserUpsertingUpdateCols   = []string{"username", "email", "role", "status", "full_name", "position", "photo",
		"notes", "security_option", "totp_secret", "password", "password_fails_in_row", "password_first_fail_at",
		"access_expire_at", "last_access", "updated_at", "deleted_at"}
	UserDefaultExcludeColumns = []string{"password", "password_fails_in_row", "password_first_fail_at"}
)

type User struct {
	ID       string `bun:",pk"`
	Username string
	Email    string `bun:",nullzero"`
	Role     base.UserRole
	Status   base.UserStatus
	FullName string
	Position string `bun:",nullzero"`
	Photo    string `bun:",nullzero"`
	Notes    string `bun:",nullzero"`

	SecurityOption      base.UserSecurityOption
	TotpSecret          string `bun:",nullzero"`
	Password            string `bun:",nullzero"`
	PasswordFailsInRow  int
	PasswordFirstFailAt time.Time `bun:",nullzero"`

	CreatedAt      time.Time `bun:",default:current_timestamp"`
	UpdatedAt      time.Time `bun:",default:current_timestamp"`
	AccessExpireAt time.Time `bun:",nullzero"`
	DeletedAt      time.Time `bun:",soft_delete,nullzero"`
	LastAccess     time.Time `bun:",nullzero"`

	Accesses []*ACLPermission `bun:"rel:has-many,join:id=subject_id"`
}

// GetID implements IDEntity interface
func (u *User) GetID() string {
	return u.ID
}

func (u *User) IsAccessExpired() bool {
	return !u.AccessExpireAt.IsZero() && u.AccessExpireAt.Before(timeutil.NowUTC())
}

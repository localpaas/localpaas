package entity

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

var (
	UserUpsertingConflictCols = []string{"id"}
	UserUpsertingUpdateCols   = []string{"email", "role", "status", "full_name", "photo", "security_option",
		"totp_secret", "password", "password_salt", "password_fails_in_row", "password_first_fail_at",
		"access_expire_at", "last_access", "updated_at", "deleted_at"}
)

type User struct {
	ID       string `bun:",pk"`
	Email    string
	Role     base.UserRole
	Status   base.UserStatus
	FullName string
	Photo    string `bun:",nullzero"`

	SecurityOption      base.UserSecurityOption
	TotpSecret          string `bun:",nullzero"`
	Password            []byte
	PasswordSalt        []byte
	PasswordFailsInRow  int
	PasswordFirstFailAt time.Time `bun:",nullzero"`

	CreatedAt      time.Time `bun:",default:current_timestamp"`
	UpdatedAt      time.Time `bun:",default:current_timestamp"`
	AccessExpireAt time.Time `bun:",nullzero"`
	DeletedAt      time.Time `bun:",soft_delete,nullzero"`
	LastAccess     time.Time `bun:",nullzero"`
}

// GetID implements IDEntity interface
func (u *User) GetID() string {
	return u.ID
}

func (u *User) IsAccessExpired() bool {
	return !u.AccessExpireAt.IsZero() && u.AccessExpireAt.Before(timeutil.NowUTC())
}

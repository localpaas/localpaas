package entity

import "time"

var (
	LoginTrustedDeviceUpsertingConflictCols = []string{"user_id", "device_id"}
	LoginTrustedDeviceUpsertingUpdateCols   = []string{"updated_at"}
)

type LoginTrustedDevice struct {
	UserID    string `bun:",pk"`
	DeviceID  string `bun:",pk"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User `bun:"rel:has-one,join:user_id=id"`
}

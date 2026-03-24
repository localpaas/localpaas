package entity

import "time"

var (
	LoginTrustedDeviceUpsertingConflictCols = []string{"user_id", "device_id"}
	LoginTrustedDeviceUpsertingUpdateCols   = []string{"updated_at"}
)

type LoginTrustedDevice struct {
	UserID    string    `bun:",pk" json:"userID"`
	DeviceID  string    `bun:",pk" json:"deviceID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User *User `bun:"rel:has-one,join:user_id=id" json:"-"`
}

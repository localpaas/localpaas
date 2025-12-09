package entity

import (
	"time"
)

type SysError struct {
	ID         string `bun:",pk"`
	RequestID  string
	Status     int
	Code       string
	Detail     string
	Cause      string
	DebugLog   string
	StackTrace string

	CreatedAt time.Time `bun:",default:current_timestamp"`
}

// GetID implements IDEntity interface
func (e *SysError) GetID() string {
	return e.ID
}

package timeutil

import "time"

type DurationMs int64

func (d DurationMs) Duration() time.Duration {
	return time.Duration(d) * time.Millisecond
}

func NewDurationMs(d time.Duration) DurationMs {
	return DurationMs(d.Milliseconds())
}

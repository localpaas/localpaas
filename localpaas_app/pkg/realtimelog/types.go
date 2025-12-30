package realtimelog

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

type LogType string

const (
	LogTypeIn   LogType = "in"
	LogTypeOut  LogType = "out"
	LogTypeErr  LogType = "err"
	LogTypeWarn LogType = "warn"
)

type LogFrame struct {
	Type LogType   `json:"type"`
	Data string    `json:"data"`
	Ts   time.Time `json:"ts"`
}

func NewInFrame(data string, ts *time.Time) *LogFrame {
	return &LogFrame{
		Type: LogTypeIn,
		Data: data,
		Ts:   calcTimestamp(ts),
	}
}

func NewOutFrame(data string, ts *time.Time) *LogFrame {
	return &LogFrame{
		Type: LogTypeOut,
		Data: data,
		Ts:   calcTimestamp(ts),
	}
}

func NewErrFrame(data string, ts *time.Time) *LogFrame {
	return &LogFrame{
		Type: LogTypeErr,
		Data: data,
		Ts:   calcTimestamp(ts),
	}
}

func NewWarnFrame(data string, ts *time.Time) *LogFrame {
	return &LogFrame{
		Type: LogTypeWarn,
		Data: data,
		Ts:   calcTimestamp(ts),
	}
}

func calcTimestamp(ts *time.Time) time.Time {
	if ts == nil {
		return timeutil.NowUTC()
	}
	return *ts
}

package timeutil

import (
	"bytes"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

type DurationStr time.Duration

func ParseDurationStr(s string) (DurationStr, error) {
	v, err := time.ParseDuration(s)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}
	return DurationStr(v), nil
}

func (dur DurationStr) ToDuration() time.Duration {
	return time.Duration(dur)
}

func (dur DurationStr) String() string {
	return time.Duration(dur).String()
}

func (dur DurationStr) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dur.String() + `"`), nil
}

func (dur *DurationStr) UnmarshalJSON(in []byte) error {
	if bytes.Equal(in, []byte("null")) {
		*dur = 0
		return nil
	}
	// Remove double quotes covering the str
	if len(in) > 1 && in[0] == '"' {
		in = in[1 : len(in)-1]
	}
	d, err := time.ParseDuration(reflectutil.UnsafeBytesToStr(in))
	if err != nil {
		return tracerr.Wrap(err)
	}
	*dur = DurationStr(d)
	return nil
}

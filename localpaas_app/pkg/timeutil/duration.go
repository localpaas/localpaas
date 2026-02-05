package timeutil

import (
	"bytes"
	"strconv"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

type Duration time.Duration

func ParseDuration(s string) (Duration, error) {
	v, err := time.ParseDuration(s)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}
	return Duration(v), nil
}

func (dur Duration) String() string {
	s := time.Duration(dur).String()
	if len(s) > 3 { //nolint:mnd
		suffix := s[len(s)-3:]
		if suffix == "m0s" || suffix == "h0m" || suffix == "d0h" {
			s = s[:len(s)-2]
		}
	}
	return s
}

func (dur Duration) ToDuration() time.Duration {
	return time.Duration(dur)
}

func (dur Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dur.String() + `"`), nil
}

func (dur *Duration) UnmarshalJSON(in []byte) error {
	if bytes.Equal(in, []byte("null")) {
		*dur = 0
		return nil
	}

	// Remove double quotes covering the str
	if len(in) > 1 && in[0] == '"' {
		in = in[1 : len(in)-1]
		d, err := time.ParseDuration(reflectutil.UnsafeBytesToStr(in))
		if err != nil {
			return tracerr.Wrap(err)
		}
		*dur = Duration(d)
		return nil
	}

	// Parse duration as integer
	v, err := strconv.ParseInt(reflectutil.UnsafeBytesToStr(in), 10, 64)
	if err != nil {
		return tracerr.Wrap(err)
	}
	*dur = Duration(v)
	return nil
}

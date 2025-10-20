package daymonth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/localpaas/localpaas/pkg/strutil"
	"github.com/localpaas/localpaas/pkg/timeutil"
	"github.com/localpaas/localpaas/pkg/tracerr"
)

type DayMonth int

var (
	ErrDayMonthFormatInvalid = errors.New("day-month format is invalid")
	ErrDayMonthValueInvalid  = errors.New("day-month value is invalid")
)

func NewFromDate(date timeutil.Date) DayMonth {
	return DayMonth(date.ToTime().Day()*100 + int(date.ToTime().Month()))
}

//nolint:mnd
func (dm DayMonth) String() string {
	return fmt.Sprintf("%02d-%02d", dm/100, dm%100)
}

//nolint:mnd
func (dm DayMonth) Valid() bool {
	dd, mm := dm/100, dm%100
	return dd >= 1 && dd <= 31 && mm >= 1 && mm <= 12
}

//nolint:mnd
func (dm DayMonth) GetMonthAndDay() (time.Month, int) {
	return time.Month(dm % 100), int(dm / 100)
}

func (dm DayMonth) GetDateForYear(year int) timeutil.Date {
	mm, dd := dm.GetMonthAndDay()
	return timeutil.Date(time.Date(year, mm, dd, 0, 0, 0, 0, time.UTC))
}

func (dm DayMonth) MarshalJSON() ([]byte, error) {
	return []byte(strutil.Quote(dm.String(), "\"")), nil
}

//nolint:mnd
func (dm *DayMonth) UnmarshalJSON(data []byte) error {
	s := strings.SplitN(strutil.Unquote(string(data), "\""), "-", 2)
	if len(s) != 2 {
		return tracerr.Wrap(ErrDayMonthFormatInvalid)
	}
	dd, err := strconv.Atoi(s[0])
	if err != nil {
		return tracerr.Wrap(err)
	}
	mm, err := strconv.Atoi(s[1])
	if err != nil {
		return tracerr.Wrap(err)
	}
	ddmm := DayMonth(dd*100 + mm)
	if !ddmm.Valid() {
		return tracerr.Wrap(ErrDayMonthValueInvalid)
	}
	*dm = ddmm
	return nil
}

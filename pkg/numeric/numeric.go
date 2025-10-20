package numeric

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/localpaas/localpaas/pkg/reflectutil"
	"github.com/localpaas/localpaas/pkg/strutil"
)

type Numeric struct {
	pgtype.Numeric
}

func (n Numeric) String() string {
	data, _ := n.Numeric.MarshalJSON()
	return reflectutil.UnsafeBytesToStr(data)
}

func (n Numeric) MarshalJSON() ([]byte, error) {
	return n.Numeric.MarshalJSON() //nolint:wrapcheck
}

func (n *Numeric) UnmarshalJSON(data []byte) error {
	dataStr := strutil.Unquote(string(data), "\"")
	return n.Numeric.UnmarshalJSON(reflectutil.UnsafeStrToBytes(dataStr)) //nolint:wrapcheck
}

func (n *Numeric) ToRat() *big.Rat {
	if n == nil {
		return nil
	}
	s := n.String()
	r := new(big.Rat)
	r.SetString(s)
	return r
}

func (n *Numeric) Cmp(other *Numeric) int {
	v1 := n.ToRat()
	v2 := other.ToRat()

	if v1 == nil && v2 == nil {
		return 0
	}

	if v1 == nil {
		return -1
	}

	return v1.Cmp(v2)
}

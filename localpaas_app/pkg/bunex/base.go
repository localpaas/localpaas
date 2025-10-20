package bunex

import "github.com/uptrace/bun"

func In(slice any) any {
	return bun.In(slice)
}

func Safe(column string) bun.Safe {
	return bun.Safe(column)
}

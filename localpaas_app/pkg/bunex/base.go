package bunex

import "github.com/uptrace/bun"

func In(slice any) any {
	return bun.In(slice)
}

func InItems(items ...any) any {
	return bun.In(items)
}

func Safe(column string) bun.Safe {
	return bun.Safe(column)
}

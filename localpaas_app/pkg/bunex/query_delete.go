package bunex

import (
	"github.com/uptrace/bun"
)

type DeleteQueryOption func(*bun.DeleteQuery) *bun.DeleteQuery

func DeleteWithDeleted() DeleteQueryOption {
	return func(query *bun.DeleteQuery) *bun.DeleteQuery {
		return query.WhereAllWithDeleted()
	}
}

func DeleteWhere(queryStr string, args ...any) DeleteQueryOption {
	return func(query *bun.DeleteQuery) *bun.DeleteQuery {
		return query.Where(queryStr, args...)
	}
}

func DeleteWhereOr(queryStr string, args ...any) DeleteQueryOption {
	return func(query *bun.DeleteQuery) *bun.DeleteQuery {
		return query.WhereOr(queryStr, args...)
	}
}

func DeleteWithForceDelete() DeleteQueryOption {
	return func(query *bun.DeleteQuery) *bun.DeleteQuery {
		return query.ForceDelete()
	}
}

// ApplyDelete applies extra delete queries to the bun query
func ApplyDelete(query *bun.DeleteQuery, opts ...DeleteQueryOption) *bun.DeleteQuery {
	for _, opt := range opts {
		query = opt(query)
	}
	return query
}

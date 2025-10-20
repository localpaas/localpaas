package bunex

import (
	"github.com/uptrace/bun"
)

type UpdateQueryOption func(*bun.UpdateQuery) *bun.UpdateQuery

func UpdateColumns(cols ...string) UpdateQueryOption {
	return func(query *bun.UpdateQuery) *bun.UpdateQuery {
		return query.Column(cols...)
	}
}

func UpdateExcludeColumns(cols ...string) UpdateQueryOption {
	return func(query *bun.UpdateQuery) *bun.UpdateQuery {
		return query.ExcludeColumn(cols...)
	}
}

func UpdateWithDeleted() UpdateQueryOption {
	return func(query *bun.UpdateQuery) *bun.UpdateQuery {
		return query.WhereAllWithDeleted()
	}
}

func UpdateWhere(queryStr string, args ...any) UpdateQueryOption {
	return func(query *bun.UpdateQuery) *bun.UpdateQuery {
		return query.Where(queryStr, args...)
	}
}

func UpdateWhereOr(queryStr string, args ...any) UpdateQueryOption {
	return func(query *bun.UpdateQuery) *bun.UpdateQuery {
		return query.WhereOr(queryStr, args...)
	}
}

// ApplyUpdate applies extra update queries to the bun query
func ApplyUpdate(query *bun.UpdateQuery, opts ...UpdateQueryOption) *bun.UpdateQuery {
	for _, opt := range opts {
		query = opt(query)
	}
	return query
}

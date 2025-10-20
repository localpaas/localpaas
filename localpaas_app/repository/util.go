package repository

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

func newPagingMeta(paging *basedto.Paging) *basedto.PagingMeta {
	if paging != nil {
		return &basedto.PagingMeta{
			Offset: paging.Offset,
			Limit:  paging.Limit,
		}
	}
	return &basedto.PagingMeta{}
}

func wrapPaginationError(err error, paging *basedto.Paging) error {
	if paging != nil && len(paging.Sort) > 0 && bunex.IsErrorColumnNotExist(err) {
		return apperrors.NewParamInvalidNT("sort").WithCause(err)
	}
	return apperrors.Wrap(err)
}

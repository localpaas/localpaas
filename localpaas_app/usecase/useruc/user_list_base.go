package useruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

func (uc *UserUC) ListUserBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.ListUserBaseReq,
) (*userdto.ListUserBaseResp, error) {
	var listOpts []bunex.SelectQueryOption

	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("\"user\".status IN (?)", bunex.In(req.Status)),
		)
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhere("\"user\".email ILIKE ?", keyword),
			bunex.SelectWhereOr("\"user\".full_name ILIKE ?", keyword),
		)
	}

	users, pagingMeta, err := uc.userRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.ListUserBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: basedto.TransformUsersBase(users),
	}, nil
}

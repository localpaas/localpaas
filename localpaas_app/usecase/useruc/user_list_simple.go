package useruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

func (uc *UserUC) ListUserSimple(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.ListUserSimpleReq,
) (*userdto.ListUserSimpleResp, error) {
	var listOpts []bunex.SelectQueryOption

	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("workspace_user.status IN (?)", bunex.In(req.Status)),
		)
	}

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			// case-insensitive search by user's full name
			bunex.SelectWhere("CONCAT(first_name, ' ', last_name) ILIKE ?", keyword),
		)
	}

	users, pagingMeta, err := uc.userRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.ListUserSimpleResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: userdto.TransformUsersSimple(users),
	}, nil
}

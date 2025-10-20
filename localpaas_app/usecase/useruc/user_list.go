package useruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

func (uc *UserUC) ListUser(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.ListUserReq,
) (*userdto.ListUserResp, error) {
	listOpts := []bunex.SelectQueryOption{
		bunex.SelectRelation("User"),
		bunex.SelectRelation("UserRoles.Role"),
		bunex.SelectRelation("Position"),
	}
	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("workspace_user.status IN (?)", bunex.In(req.Status)),
		)
	}
	// Filter by search keyword
	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				// NOTE: Needs to quote `user` as it is keyword
				bunex.SelectWhere("\"user\".email ILIKE ?", keyword),
				bunex.SelectWhereOr("CONCAT(\"user\".first_name, ' ', \"user\".last_name) ILIKE ?",
					keyword),
			),
		)
	}

	users, paging, err := uc.userRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := userdto.TransformUsers(users)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.ListUserResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}

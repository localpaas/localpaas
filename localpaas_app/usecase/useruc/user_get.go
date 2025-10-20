package useruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

func (uc *UserUC) GetUser(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.GetUserReq,
) (*userdto.GetUserResp, error) {
	user, err := uc.userRepo.GetByID(ctx, uc.db, req.ID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := userdto.TransformUserDetails(user)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.GetUserResp{
		Data: resp,
	}, nil
}

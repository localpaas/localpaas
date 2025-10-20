package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) GetMe(
	_ context.Context,
	auth *basedto.Auth,
	_ *sessiondto.GetMeReq,
) (*sessiondto.GetMeResp, error) {
	userResp, err := sessiondto.TransformUser(auth.User.User)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sessiondto.GetMeResp{
		Data: &sessiondto.GetMeDataResp{
			User: userResp,
		},
	}, nil
}

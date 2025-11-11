package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) GetMe(
	_ context.Context,
	user *basedto.User,
	_ *sessiondto.GetMeReq,
) (*sessiondto.GetMeResp, error) {
	userResp, err := sessiondto.TransformUser(user.User)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp := &sessiondto.GetMeDataResp{
		User: userResp,
	}

	if user.Status == base.UserStatusPending && user.TotpSecret == "" {
		resp.NextStep = nextStepMfaSetup
	}

	return &sessiondto.GetMeResp{
		Data: resp,
	}, nil
}

package sessionuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

func (uc *SessionUC) GetLoginOptions(
	ctx context.Context,
	req *sessiondto.GetLoginOptionsReq,
) (resp *sessiondto.GetLoginOptionsResp, err error) {
	// TODO: implement this
	return &sessiondto.GetLoginOptionsResp{
		Data: &sessiondto.LoginOptionsResp{
			AllowLoginWithGitHub: false,
			AllowLoginWithGitLab: false,
		},
	}, nil
}

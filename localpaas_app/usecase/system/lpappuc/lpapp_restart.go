package lpappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/lpappuc/lpappdto"
)

func (uc *LpAppUC) RestartLpApp(
	ctx context.Context,
	_ *basedto.Auth,
	_ *lpappdto.RestartLpAppReq,
) (*lpappdto.RestartLpAppResp, error) {
	err := uc.lpAppService.RestartLpAppSwarmService(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &lpappdto.RestartLpAppResp{}, nil
}

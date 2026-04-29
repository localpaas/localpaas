package lpappuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/lpappuc/lpappdto"
)

func (uc *UC) GetLpAppReleaseInfo(
	ctx context.Context,
	_ *basedto.Auth,
	_ *lpappdto.GetLpAppReleaseInfoReq,
) (*lpappdto.GetLpAppReleaseInfoResp, error) {
	info, err := uc.lpAppService.GetAppReleaseInfo(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &lpappdto.GetLpAppReleaseInfoResp{
		Data: &lpappdto.LpAppReleaseInfoResp{
			AppReleaseInfo: info,
		},
	}, nil
}

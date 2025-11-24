package nginxuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc/nginxdto"
)

func (uc *NginxUC) ResetNginxConfig(
	ctx context.Context,
	_ *basedto.Auth,
	_ *nginxdto.ResetNginxConfigReq,
) (*nginxdto.ResetNginxConfigResp, error) {
	err := uc.nginxService.ResetNginxConfig(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &nginxdto.ResetNginxConfigResp{}, nil
}

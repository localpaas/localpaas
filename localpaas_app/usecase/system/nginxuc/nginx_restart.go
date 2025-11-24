package nginxuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/system/nginxuc/nginxdto"
)

func (uc *NginxUC) RestartNginx(
	ctx context.Context,
	_ *basedto.Auth,
	_ *nginxdto.RestartNginxReq,
) (*nginxdto.RestartNginxResp, error) {
	err := uc.nginxService.RestartNginxSwarmService(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &nginxdto.RestartNginxResp{}, nil
}

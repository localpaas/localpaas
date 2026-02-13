package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

func (uc *SSLUC) UpdateSSLMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.UpdateSSLMetaReq,
) (*ssldto.UpdateSSLMetaResp, error) {
	req.Type = currentSettingType
	_, err := uc.UpdateSettingMeta(ctx, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.UpdateSSLMetaResp{}, nil
}

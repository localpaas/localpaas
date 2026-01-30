package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

func (uc *SSLUC) ListSSL(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.ListSSLReq,
) (*ssldto.ListSSLResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := ssldto.TransformSSLs(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.ListSSLResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

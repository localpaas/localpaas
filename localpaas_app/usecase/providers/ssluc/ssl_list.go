package ssluc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/ssluc/ssldto"
)

func (uc *SslUC) ListSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *ssldto.ListSslReq,
) (*ssldto.ListSslResp, error) {
	req.Type = currentSettingType
	resp, err := providers.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &providers.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := ssldto.TransformSsls(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &ssldto.ListSslResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

package sslcertsettingsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sslcertsettingsuc/sslcertsettingsdto"
)

func (uc *UC) GetUniqueSSLCertSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *sslcertsettingsdto.GetUniqueSSLCertSettingsReq,
) (*sslcertsettingsdto.GetUniqueSSLCertSettingsResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetUniqueSetting(ctx, auth, &req.GetUniqueSettingReq, &settings.GetUniqueSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := sslcertsettingsdto.TransformSSLCertSettings(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sslcertsettingsdto.GetUniqueSSLCertSettingsResp{
		Data: respData,
	}, nil
}

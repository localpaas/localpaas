package sslcertsettingsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sslcertsettingsuc/sslcertsettingsdto"
)

func (uc *UC) UpdateUniqueSSLCertSettingsMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaReq,
) (*sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaResp, error) {
	req.Type = currentSettingType
	_, err := uc.UpdateUniqueSettingMeta(ctx, &req.UpdateUniqueSettingMetaReq, &settings.UpdateUniqueSettingMetaData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sslcertsettingsdto.UpdateUniqueSSLCertSettingsMetaResp{}, nil
}

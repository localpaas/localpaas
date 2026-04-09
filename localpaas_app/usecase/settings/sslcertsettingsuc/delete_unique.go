package sslcertsettingsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sslcertsettingsuc/sslcertsettingsdto"
)

func (uc *UC) DeleteUniqueSSLCertSettings(
	ctx context.Context,
	auth *basedto.Auth,
	req *sslcertsettingsdto.DeleteUniqueSSLCertSettingsReq,
) (*sslcertsettingsdto.DeleteUniqueSSLCertSettingsResp, error) {
	req.Type = currentSettingType
	_, err := uc.DeleteUniqueSetting(ctx, &req.DeleteUniqueSettingReq, &settings.DeleteUniqueSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sslcertsettingsdto.DeleteUniqueSSLCertSettingsResp{}, nil
}

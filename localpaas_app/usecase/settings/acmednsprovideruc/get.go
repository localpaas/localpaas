package acmednsprovideruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/acmednsprovideruc/acmednsproviderdto"
)

func (uc *UC) GetAcmeDnsProvider(
	ctx context.Context,
	auth *basedto.Auth,
	req *acmednsproviderdto.GetAcmeDnsProviderReq,
) (*acmednsproviderdto.GetAcmeDnsProviderResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting := resp.Data
	if setting.ObjectID == setting.CurrentObjectID { // not return sensitive data if setting is inherited
		setting.MustAsAcmeDnsProvider().MustDecrypt()
	}

	respData, err := acmednsproviderdto.TransformAcmeDnsProvider(setting, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &acmednsproviderdto.GetAcmeDnsProviderResp{
		Data: respData,
	}, nil
}

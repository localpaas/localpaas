package cloudprovideruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
)

func (uc *CloudProviderUC) GetCloudProvider(
	ctx context.Context,
	auth *basedto.Auth,
	req *cloudproviderdto.GetCloudProviderReq,
) (*cloudproviderdto.GetCloudProviderResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Data.MustAsCloudProvider().MustDecrypt()
	respData, err := cloudproviderdto.TransformCloudProvider(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cloudproviderdto.GetCloudProviderResp{
		Data: respData,
	}, nil
}

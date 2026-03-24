package cloudprovideruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
)

func (uc *CloudProviderUC) DeleteCloudProvider(
	ctx context.Context,
	auth *basedto.Auth,
	req *cloudproviderdto.DeleteCloudProviderReq,
) (*cloudproviderdto.DeleteCloudProviderResp, error) {
	req.Type = currentSettingType
	_, err := uc.DeleteSetting(ctx, &req.DeleteSettingReq, &settings.DeleteSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cloudproviderdto.DeleteCloudProviderResp{}, nil
}

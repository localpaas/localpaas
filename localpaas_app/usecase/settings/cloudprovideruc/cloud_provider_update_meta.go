package cloudprovideruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
)

func (uc *CloudProviderUC) UpdateCloudProviderMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *cloudproviderdto.UpdateCloudProviderMetaReq,
) (*cloudproviderdto.UpdateCloudProviderMetaResp, error) {
	req.Type = currentSettingType
	_, err := uc.UpdateSettingMeta(ctx, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cloudproviderdto.UpdateCloudProviderMetaResp{}, nil
}

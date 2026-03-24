package cloudprovideruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
)

func (uc *CloudProviderUC) ListCloudProvider(
	ctx context.Context,
	auth *basedto.Auth,
	req *cloudproviderdto.ListCloudProviderReq,
) (*cloudproviderdto.ListCloudProviderResp, error) {
	req.Type = currentSettingType
	resp, err := uc.ListSetting(ctx, auth, &req.ListSettingReq, &settings.ListSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := cloudproviderdto.TransformAWSs(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &cloudproviderdto.ListCloudProviderResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

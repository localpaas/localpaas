package awsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
)

func (uc *AWSUC) UpdateAWSMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *awsdto.UpdateAWSMetaReq,
) (*awsdto.UpdateAWSMetaResp, error) {
	req.Type = currentSettingType
	_, err := uc.UpdateSettingMeta(ctx, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awsdto.UpdateAWSMetaResp{}, nil
}

package awss3uc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
)

func (uc *AWSS3UC) DeleteAWSS3(
	ctx context.Context,
	auth *basedto.Auth,
	req *awss3dto.DeleteAWSS3Req,
) (*awss3dto.DeleteAWSS3Resp, error) {
	req.Type = currentSettingType
	_, err := uc.DeleteSetting(ctx, &req.DeleteSettingReq, &settings.DeleteSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awss3dto.DeleteAWSS3Resp{}, nil
}

package awss3uc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
)

func (uc *AWSS3UC) GetAWSS3(
	ctx context.Context,
	auth *basedto.Auth,
	req *awss3dto.GetAWSS3Req,
) (*awss3dto.GetAWSS3Resp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Data.MustAsAWSS3().MustDecrypt()
	respData, err := awss3dto.TransformAWSS3(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awss3dto.GetAWSS3Resp{
		Data: respData,
	}, nil
}

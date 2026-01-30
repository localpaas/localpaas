package awss3uc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
)

func (uc *AWSS3UC) ListAWSS3(
	ctx context.Context,
	auth *basedto.Auth,
	req *awss3dto.ListAWSS3Req,
) (*awss3dto.ListAWSS3Resp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := awss3dto.TransformAWSS3s(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awss3dto.ListAWSS3Resp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

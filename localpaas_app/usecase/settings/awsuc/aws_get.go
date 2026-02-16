package awsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
)

func (uc *AWSUC) GetAWS(
	ctx context.Context,
	auth *basedto.Auth,
	req *awsdto.GetAWSReq,
) (*awsdto.GetAWSResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp.Data.MustAsAWS().MustDecrypt()
	respData, err := awsdto.TransformAWS(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awsdto.GetAWSResp{
		Data: respData,
	}, nil
}

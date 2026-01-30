package awsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
)

func (uc *AWSUC) ListAWS(
	ctx context.Context,
	auth *basedto.Auth,
	req *awsdto.ListAWSReq,
) (*awsdto.ListAWSResp, error) {
	req.Type = currentSettingType
	resp, err := settings.ListSetting(ctx, uc.db, auth, &req.ListSettingReq, &settings.ListSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := awsdto.TransformAWSs(resp.Data)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &awsdto.ListAWSResp{
		Meta: resp.Meta,
		Data: respData,
	}, nil
}

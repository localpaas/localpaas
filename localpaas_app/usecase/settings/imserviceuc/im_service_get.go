package imserviceuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
)

func (uc *IMServiceUC) GetIMService(
	ctx context.Context,
	auth *basedto.Auth,
	req *imservicedto.GetIMServiceReq,
) (*imservicedto.GetIMServiceResp, error) {
	req.Type = currentSettingType
	setting, err := uc.GetSetting(ctx, auth, &req.GetSettingReq, &settings.GetSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsIMService().MustDecrypt()
	resp, err := imservicedto.TransformIMService(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imservicedto.GetIMServiceResp{
		Data: resp,
	}, nil
}

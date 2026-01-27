package emailuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
)

func (uc *EmailUC) GetEmail(
	ctx context.Context,
	auth *basedto.Auth,
	req *emaildto.GetEmailReq,
) (*emaildto.GetEmailResp, error) {
	req.Type = currentSettingType
	setting, err := settings.GetSetting(ctx, uc.db, auth, &req.GetSettingReq, &settings.GetSettingData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	setting.MustAsEmail().MustDecrypt()
	resp, err := emaildto.TransformEmail(setting)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &emaildto.GetEmailResp{
		Data: resp,
	}, nil
}

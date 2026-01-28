package imserviceuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
)

func (uc *IMServiceUC) UpdateIMServiceMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *imservicedto.UpdateIMServiceMetaReq,
) (*imservicedto.UpdateIMServiceMetaResp, error) {
	req.Type = currentSettingType
	_, err := settings.UpdateSettingMeta(ctx, uc.db, &req.UpdateSettingMetaReq, &settings.UpdateSettingMetaData{
		SettingRepo: uc.settingRepo,
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imservicedto.UpdateIMServiceMetaResp{}, nil
}

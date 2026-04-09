package imagebuildsettingsuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuildsettingsuc/imagebuildsettingsdto"
)

func (uc *UC) UpdateUniqueImageBuildSettingsMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaReq,
) (*imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaResp, error) {
	req.Type = currentSettingType
	_, err := uc.UpdateUniqueSettingMeta(ctx, &req.UpdateUniqueSettingMetaReq, &settings.UpdateUniqueSettingMetaData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imagebuildsettingsdto.UpdateUniqueImageBuildSettingsMetaResp{}, nil
}

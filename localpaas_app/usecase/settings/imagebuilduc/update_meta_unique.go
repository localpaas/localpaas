package imagebuilduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

func (uc *UC) UpdateUniqueImageBuildMeta(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagebuilddto.UpdateUniqueImageBuildMetaReq,
) (*imagebuilddto.UpdateUniqueImageBuildMetaResp, error) {
	req.Type = currentSettingType
	_, err := uc.UpdateUniqueSettingMeta(ctx, &req.UpdateUniqueSettingMetaReq, &settings.UpdateUniqueSettingMetaData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imagebuilddto.UpdateUniqueImageBuildMetaResp{}, nil
}

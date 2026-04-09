package imagebuilduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

func (uc *UC) DeleteUniqueImageBuild(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagebuilddto.DeleteUniqueImageBuildReq,
) (*imagebuilddto.DeleteUniqueImageBuildResp, error) {
	req.Type = currentSettingType
	_, err := uc.DeleteUniqueSetting(ctx, &req.DeleteUniqueSettingReq, &settings.DeleteUniqueSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imagebuilddto.DeleteUniqueImageBuildResp{}, nil
}

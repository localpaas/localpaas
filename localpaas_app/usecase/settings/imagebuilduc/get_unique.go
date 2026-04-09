package imagebuilduc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imagebuilduc/imagebuilddto"
)

func (uc *UC) GetUniqueImageBuild(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagebuilddto.GetUniqueImageBuildReq,
) (*imagebuilddto.GetUniqueImageBuildResp, error) {
	req.Type = currentSettingType
	resp, err := uc.GetUniqueSetting(ctx, auth, &req.GetUniqueSettingReq, &settings.GetUniqueSettingData{})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	respData, err := imagebuilddto.TransformImageBuild(resp.Data, resp.RefObjects)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &imagebuilddto.GetUniqueImageBuildResp{
		Data: respData,
	}, nil
}

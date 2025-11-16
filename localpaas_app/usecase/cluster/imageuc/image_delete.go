package imageuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/imageuc/imagedto"
)

func (uc *ImageUC) DeleteImage(
	ctx context.Context,
	auth *basedto.Auth,
	req *imagedto.DeleteImageReq,
) (*imagedto.DeleteImageResp, error) {
	_, err := uc.dockerManager.ImageRemove(ctx, req.ImageID)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}

	return &imagedto.DeleteImageResp{}, nil
}

package volumeuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc/volumedto"
)

func (uc *VolumeUC) GetVolume(
	ctx context.Context,
	auth *basedto.Auth,
	req *volumedto.GetVolumeReq,
) (*volumedto.GetVolumeResp, error) {
	vol, _, err := uc.dockerManager.VolumeInspect(ctx, req.VolumeID)
	if err != nil {
		return nil, apperrors.NewInfra(err)
	}

	return &volumedto.GetVolumeResp{
		Data: volumedto.TransformVolume(vol, true),
	}, nil
}

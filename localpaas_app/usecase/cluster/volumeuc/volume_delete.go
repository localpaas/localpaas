package volumeuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc/volumedto"
)

func (uc *VolumeUC) DeleteVolume(
	ctx context.Context,
	auth *basedto.Auth,
	req *volumedto.DeleteVolumeReq,
) (*volumedto.DeleteVolumeResp, error) {
	err := uc.dockerManager.VolumeRemove(ctx, req.VolumeID, req.Force)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &volumedto.DeleteVolumeResp{}, nil
}

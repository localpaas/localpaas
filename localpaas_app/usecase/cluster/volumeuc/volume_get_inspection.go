package volumeuc

import (
	"context"
	"encoding/json"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/volumeuc/volumedto"
)

func (uc *VolumeUC) GetVolumeInspection(
	ctx context.Context,
	auth *basedto.Auth,
	req *volumedto.GetVolumeInspectionReq,
) (*volumedto.GetVolumeInspectionResp, error) {
	vol, _, err := uc.dockerManager.VolumeInspect(ctx, req.VolumeID)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := json.MarshalIndent(vol, "", "   ")
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &volumedto.GetVolumeInspectionResp{
		Data: reflectutil.UnsafeBytesToStr(resp),
	}, nil
}

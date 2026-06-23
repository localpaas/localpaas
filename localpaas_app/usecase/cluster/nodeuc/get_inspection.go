package nodeuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/nodeuc/nodedto"
)

func (uc *UC) GetNodeInspection(
	ctx context.Context,
	auth *basedto.Auth,
	req *nodedto.GetNodeInspectionReq,
) (*nodedto.GetNodeInspectionResp, error) {
	resp, err := uc.dockerManager.NodeInspect(ctx, req.NodeID)
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &nodedto.GetNodeInspectionResp{
		Data: reflectutil.UnsafeBytesToStr(resp.Raw),
	}, nil
}

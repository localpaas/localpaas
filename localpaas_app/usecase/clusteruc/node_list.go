package clusteruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) ListNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.ListNodeReq,
) (*clusterdto.ListNodeResp, error) {
	listOpts := []bunex.SelectQueryOption{}
	if len(req.Status) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("node.status IN (?)", bunex.In(req.Status)),
		)
	}
	if len(req.InfraStatus) > 0 {
		listOpts = append(listOpts,
			bunex.SelectWhere("node.infra_status IN (?)", bunex.In(req.InfraStatus)),
		)
	}
	// Filter by search keyword
	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("node.host_name ILIKE ?", keyword),
				bunex.SelectWhereOr("node.ip ILIKE ?", keyword),
				bunex.SelectWhereOr("node.note ILIKE ?", keyword),
			),
		)
	}

	nodes, paging, err := uc.nodeRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	resp, err := clusterdto.TransformNodes(nodes)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &clusterdto.ListNodeResp{
		Meta: &basedto.Meta{Page: paging},
		Data: resp,
	}, nil
}

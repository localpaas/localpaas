package clusteruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) ListNodeBase(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.ListNodeBaseReq,
) (*clusterdto.ListNodeBaseResp, error) {
	var listOpts []bunex.SelectQueryOption

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

	if req.Search != "" {
		keyword := bunex.MakeLikeOpStr(req.Search, true)
		listOpts = append(listOpts,
			bunex.SelectWhereGroup(
				bunex.SelectWhere("node.name ILIKE ?", keyword),
				bunex.SelectWhereOr("node.ip ILIKE ?", keyword),
				bunex.SelectWhereOr("node.note ILIKE ?", keyword),
			),
		)
	}

	nodes, pagingMeta, err := uc.nodeRepo.List(ctx, uc.db, &req.Paging, listOpts...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &clusterdto.ListNodeBaseResp{
		Meta: &basedto.Meta{Page: pagingMeta},
		Data: clusterdto.TransformNodesBase(nodes),
	}, nil
}

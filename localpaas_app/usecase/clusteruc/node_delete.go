package clusteruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/service/clusterservice"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

func (uc *ClusterUC) DeleteNode(
	ctx context.Context,
	auth *basedto.Auth,
	req *clusterdto.DeleteNodeReq,
) (*clusterdto.DeleteNodeResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		appData := &deleteNodeData{}
		err := uc.loadNodeDataForDelete(ctx, db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &clusterservice.PersistingClusterData{}
		uc.prepareDeletingNode(appData, persistingData)

		// TODO: handle node deletion

		return uc.clusterService.PersistClusterData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &clusterdto.DeleteNodeResp{}, nil
}

type deleteNodeData struct {
	Node *entity.Node
}

func (uc *ClusterUC) loadNodeDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *clusterdto.DeleteNodeReq,
	data *deleteNodeData,
) error {
	node, err := uc.nodeRepo.GetByID(ctx, db, req.NodeID,
		bunex.SelectFor("UPDATE OF node"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Node = node

	if node.Status == base.NodeStatusDeleting { //nolint
		// TODO: handle task deletion if previously failed
	}

	return nil
}

func (uc *ClusterUC) prepareDeletingNode(
	data *deleteNodeData,
	persistingData *clusterservice.PersistingClusterData,
) {
	app := data.Node
	app.Status = base.NodeStatusDeleting
	persistingData.UpsertingNodes = append(persistingData.UpsertingNodes, app)
}

package clusterhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/clusteruc/clusterdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListNode Lists cluster nodes
// @Summary Lists cluster nodes
// @Description Lists cluster nodes
// @Tags    cluster_nodes
// @Produce json
// @Id      listNode
// @Param   status query string false "`status=<target>`"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} clusterdto.ListNodeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/nodes [get]
func (h *ClusterHandler) ListNode(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		RequireAdmin: true,
		ResourceType: base.ResourceTypeNode,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := clusterdto.NewListNodeReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.clusterUC.ListNode(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetNode Gets node details
// @Summary Gets node details
// @Description Gets node details
// @Tags    cluster_nodes
// @Produce json
// @Id      getNode
// @Param   nodeID path string true "node ID"
// @Success 200 {object} clusterdto.GetNodeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/nodes/{nodeID} [get]
func (h *ClusterHandler) GetNode(ctx *gin.Context) {
	nodeID, err := h.ParseStringParam(ctx, "nodeID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		RequireAdmin: true,
		ResourceType: base.ResourceTypeNode,
		ResourceID:   nodeID,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := clusterdto.NewGetNodeReq()
	req.NodeID = nodeID
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.clusterUC.GetNode(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateNode Updates a node
// @Summary Updates a node
// @Description Updates a node
// @Tags    cluster_nodes
// @Produce json
// @Id      updateNode
// @Param   nodeID path string true "node ID"
// @Param   body body clusterdto.UpdateNodeReq true "request data"
// @Success 200 {object} clusterdto.UpdateNodeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/nodes/{nodeID} [put]
func (h *ClusterHandler) UpdateNode(ctx *gin.Context) {
	nodeID, err := h.ParseStringParam(ctx, "nodeID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		RequireAdmin: true,
		ResourceType: base.ResourceTypeNode,
		ResourceID:   nodeID,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := clusterdto.NewUpdateNodeReq()
	req.NodeID = nodeID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.clusterUC.UpdateNode(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteNode Deletes a node
// @Summary Deletes a node
// @Description Deletes a node
// @Tags    cluster_nodes
// @Produce json
// @Id      deleteNode
// @Param   nodeID path string true "node ID"
// @Success 200 {object} clusterdto.DeleteNodeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/nodes/{nodeID} [delete]
func (h *ClusterHandler) DeleteNode(ctx *gin.Context) {
	nodeID, err := h.ParseStringParam(ctx, "nodeID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		RequireAdmin: true,
		ResourceType: base.ResourceTypeNode,
		ResourceID:   nodeID,
		Action:       base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := clusterdto.NewDeleteNodeReq()
	req.NodeID = nodeID
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.clusterUC.DeleteNode(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

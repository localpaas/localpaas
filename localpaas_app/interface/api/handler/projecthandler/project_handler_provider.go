package projecthandler

import (
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

func (h *ProjectHandler) getProjectProviderAuth(
	ctx *gin.Context,
	action base.ActionType,
	getProviderItemId bool,
) (auth *basedto.Auth, projectID string, itemID string, err error) {
	projectID, err = h.ParseStringParam(ctx, "projectID")
	if err != nil {
		return
	}
	if getProviderItemId {
		itemID, err = h.ParseStringParam(ctx, "id")
		if err != nil {
			return
		}
	}
	auth, err = h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProject,
		ResourceType:   base.ResourceTypeProject,
		ResourceID:     projectID,
		Action:         action,
	})
	if err != nil {
		return
	}
	return
}

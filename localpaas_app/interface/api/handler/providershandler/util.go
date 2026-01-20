package providershandler

import (
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

func (h *ProvidersHandler) getAuth(
	ctx *gin.Context,
	resourceType base.ResourceType,
	action base.ActionType,
	getItemID bool,
) (auth *basedto.Auth, itemID string, err error) {
	if getItemID {
		itemID, err = h.ParseStringParam(ctx, "id")
		if err != nil {
			return
		}
	}
	auth, err = h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   resourceType,
		ResourceID:     itemID,
		Action:         action,
	})
	if err != nil {
		return
	}
	return
}

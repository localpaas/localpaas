package apphandler

import (
	"github.com/gin-gonic/gin"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

//nolint:nakedret
func (h *AppHandler) getAuth(
	ctx *gin.Context,
	action base.ActionType,
	getAppID bool,
) (auth *basedto.Auth, projectID, appID string, err error) {
	projectID, err = h.ParseStringParam(ctx, "projectID")
	if err != nil {
		return
	}
	var accessCheck *permission.AccessCheck
	if getAppID {
		appID, err = h.ParseStringParam(ctx, "appID")
		if err != nil {
			return
		}
		accessCheck = &permission.AccessCheck{
			ResourceModule:     base.ResourceModuleProject,
			ResourceType:       base.ResourceTypeApp,
			ResourceID:         appID,
			ParentResourceType: base.ResourceTypeProject,
			ParentResourceID:   projectID,
			Action:             action,
		}
	} else {
		accessCheck = &permission.AccessCheck{
			ResourceModule: base.ResourceModuleProject,
			ResourceType:   base.ResourceTypeProject,
			ResourceID:     projectID,
			Action:         gofn.If(action == base.ActionTypeDelete, base.ActionTypeWrite, action),
		}
	}
	auth, err = h.authHandler.GetCurrentAuth(ctx, accessCheck)
	if err != nil {
		return
	}
	return
}

//nolint:nakedret
func (h *AppHandler) getAuthForItem(
	ctx *gin.Context,
	action base.ActionType,
) (auth *basedto.Auth, projectID, appID, itemID string, err error) {
	projectID, err = h.ParseStringParam(ctx, "projectID")
	if err != nil {
		return
	}
	appID, err = h.ParseStringParam(ctx, "appID")
	if err != nil {
		return
	}
	itemID, err = h.ParseStringParam(ctx, "id")
	if err != nil {
		return
	}
	auth, err = h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		Action:             action,
	})
	if err != nil {
		return
	}
	return
}

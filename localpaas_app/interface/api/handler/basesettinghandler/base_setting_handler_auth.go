package basesettinghandler

import (
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

func (h *BaseSettingHandler) GetAuthGlobalSettings(
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
	auth, err = h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   resourceType,
		ResourceID:     itemID,
		Action:         action,
	})
	if err != nil {
		return
	}
	return
}

func (h *BaseSettingHandler) GetAuthUserSettings(
	ctx *gin.Context,
	_ base.ActionType,
	getItemID bool,
) (auth *basedto.Auth, itemID string, err error) {
	auth, err = h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		return
	}
	if getItemID {
		itemID, err = h.ParseStringParam(ctx, "id")
		if err != nil {
			return
		}
	}
	return
}

func (h *BaseSettingHandler) GetAuthProjectSettings(
	ctx *gin.Context,
	action base.ActionType,
	getItemID bool,
) (auth *basedto.Auth, projectID, itemID string, err error) {
	projectID, err = h.ParseStringParam(ctx, "projectID")
	if err != nil {
		return
	}
	if getItemID {
		itemID, err = h.ParseStringParam(ctx, "id")
		if err != nil {
			return
		}
	}
	auth, err = h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
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

//nolint:nakedret
func (h *BaseSettingHandler) GetAuthAppSettings(
	ctx *gin.Context,
	action base.ActionType,
	getItemID bool,
) (auth *basedto.Auth, projectID, appID, itemID string, err error) {
	projectID, err = h.ParseStringParam(ctx, "projectID")
	if err != nil {
		return
	}
	appID, err = h.ParseStringParam(ctx, "appID")
	if err != nil {
		return
	}
	if getItemID {
		itemID, err = h.ParseStringParam(ctx, "id")
		if err != nil {
			return
		}
	}
	auth, err = h.AuthHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             action,
	})
	if err != nil {
		return
	}
	return
}

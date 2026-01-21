package userhandler

import (
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/permission"
)

func (h *UserHandler) getAuth(
	ctx *gin.Context,
	resType base.ResourceType,
	action base.ActionType,
	getUserID bool,
) (auth *basedto.Auth, userID string, err error) {
	if getUserID {
		userID, err = h.ParseStringParam(ctx, "userID")
		if err != nil {
			return
		}
	}
	auth, err = h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleUser,
		ResourceType:   resType,
		ResourceID:     userID,
		Action:         action,
	})
	if err != nil {
		return
	}
	return
}

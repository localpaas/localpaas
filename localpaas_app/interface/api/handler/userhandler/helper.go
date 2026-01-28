package userhandler

import (
	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
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
	accessCheck := &permission.AccessCheck{
		ResourceModule: base.ResourceModuleUser,
		ResourceType:   resType,
		ResourceID:     userID,
		Action:         action,
	}
	if userID == "current" {
		accessCheck = authhandler.NoAccessCheck
	}
	auth, err = h.authHandler.GetCurrentAuth(ctx, accessCheck)
	if auth != nil && (userID == "current" || userID == auth.User.ID) {
		err = nil
		userID = auth.User.ID
	}
	if err != nil {
		return
	}
	return
}

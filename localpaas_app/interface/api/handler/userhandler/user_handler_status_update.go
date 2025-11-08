package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// UpdateUserStatus Updates user status
// @Summary Updates user status
// @Description Updates user status
// @Tags    users
// @Produce json
// @Id      updateUserStatus
// @Param   userID path string true "user ID"
// @Param   body body userdto.UpdateStatusReq true "request data"
// @Success 200 {object} userdto.UpdateStatusResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/{userID}/status [put]
func (h *UserHandler) UpdateUserStatus(ctx *gin.Context) {
	userID, err := h.ParseStringParam(ctx, "userID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		RequireAdmin: true,
		ResourceType: base.ResourceTypeUser,
		ResourceID:   userID,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewUpdateStatusReq()
	req.ID = userID
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.UpdateStatus(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

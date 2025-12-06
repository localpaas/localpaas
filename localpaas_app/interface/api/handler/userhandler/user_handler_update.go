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

// UpdateUser Updates user data (admin API)
// @Summary Updates user data (admin API)
// @Description Updates user data (admin API)
// @Tags    users
// @Produce json
// @Id      updateUser
// @Param   userID path string true "user ID"
// @Param   body body userdto.UpdateUserReq true "request data"
// @Success 200 {object} userdto.UpdateUserResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/{userID} [put]
func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	userID, err := h.ParseStringParam(ctx, "userID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleUser,
		ResourceType:   base.ResourceTypeUser,
		ResourceID:     userID,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewUpdateUserReq()
	req.ID = userID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.UpdateUser(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteUser Deletes a user (admin API)
// @Summary Deletes a user (admin API)
// @Description Deletes a user (admin API)
// @Tags    users
// @Produce json
// @Id      deleteUser
// @Param   userID path string true "user ID"
// @Success 200 {object} userdto.DeleteUserResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/{userID} [delete]
func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userID, err := h.ParseStringParam(ctx, "userID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleUser,
		ResourceType:   base.ResourceTypeUser,
		ResourceID:     userID,
		Action:         base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewDeleteUserReq()
	req.ID = userID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.DeleteUser(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

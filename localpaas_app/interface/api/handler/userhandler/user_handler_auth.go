package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// UpdateUserPassword Changes user password
// @Summary Changes user password
// @Description Changes user password
// @Tags    users
// @Produce json
// @Id      updateUserPassword
// @Param   body body userdto.UpdatePasswordReq true "request data"
// @Success 200 {object} userdto.UpdatePasswordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/current/password [put]
func (h *UserHandler) UpdateUserPassword(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewUpdatePasswordReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.UpdatePassword(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RequestResetPassword Requests password reset
// @Summary Requests password reset
// @Description Requests password reset
// @Tags    users
// @Produce json
// @Id      requestResetPassword
// @Param   userID path string true "user ID"
// @Param   body body userdto.RequestResetPasswordReq true "request data"
// @Success 200 {object} userdto.RequestResetPasswordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/{userID}/password/request-reset [post]
func (h *UserHandler) RequestResetPassword(ctx *gin.Context) {
	userID, err := h.ParseStringParam(ctx, "userID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleUser,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewRequestResetPasswordReq()
	req.ID = userID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.RequestResetPassword(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ResetPassword Resets password
// @Summary Resets password
// @Description Resets password
// @Tags    users
// @Produce json
// @Id      resetPassword
// @Param   userID path string true "user ID"
// @Param   body body userdto.ResetPasswordReq true "request data"
// @Success 200 {object} userdto.ResetPasswordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/{userID}/password/reset [post]
func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	userID, err := h.ParseStringParam(ctx, "userID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewResetPasswordReq()
	req.ID = userID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.ResetPassword(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

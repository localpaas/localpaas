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

// InviteUser Invites a user to system
// @Summary Invites a user to system
// @Description Invites a user to system
// @Tags    users
// @Produce json
// @Id      inviteUser
// @Param   body body userdto.InviteUserReq true "request data"
// @Success 200 {object} userdto.InviteUserResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/invite [post]
func (h *UserHandler) InviteUser(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleUser,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewInviteUserReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.InviteUser(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// BeginUserSignup Begins user signup process
// @Summary Begins user signup process
// @Description Begins user signup process
// @Tags    users
// @Produce json
// @Id      beginUserSignup
// @Param   body body userdto.BeginUserSignupReq true "request data"
// @Success 200 {object} userdto.BeginUserSignupResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/signup-begin [post]
func (h *UserHandler) BeginUserSignup(ctx *gin.Context) {
	req := userdto.NewBeginUserSignupReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.BeginUserSignup(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CompleteUserSignup Completes user signup process
// @Summary Completes user signup process
// @Description Completes user signup process
// @Tags    users
// @Produce json
// @Id      completeUserSignup
// @Param   body body userdto.CompleteUserSignupReq true "request data"
// @Success 200 {object} userdto.CompleteUserSignupResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/signup-complete [post]
func (h *UserHandler) CompleteUserSignup(ctx *gin.Context) {
	req := userdto.NewCompleteUserSignupReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.userUC.CompleteUserSignup(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

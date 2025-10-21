package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// UpdateUserPassword Changes user password
// @Summary Changes user password
// @Description Changes user password
// @Tags    users_auth
// @Produce json
// @Id      updateUserPassword
// @Param   body body userdto.UpdatePasswordReq true "request data"
// @Success 200 {object} userdto.UpdatePasswordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/current/password [patch]
func (h *UserHandler) UpdateUserPassword(ctx *gin.Context) {
	// Every user can change their own password
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := userdto.NewUpdatePasswordReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
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

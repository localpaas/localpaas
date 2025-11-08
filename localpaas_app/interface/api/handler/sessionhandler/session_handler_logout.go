package sessionhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// DeleteSession Deletes the current user session
// @Summary Deletes the current user session
// @Description Deletes the current user session
// @Tags    sessions
// @Produce json
// @Id      deleteSession
// @Success 200 {object} sessiondto.DeleteSessionResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /sessions [delete]
func (h *SessionHandler) DeleteSession(ctx *gin.Context) {
	user, err := h.authHandler.GetCurrentUser(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sessiondto.NewDeleteSessionReq()
	req.User = user
	if err = h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sessionUC.DeleteSession(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Clear session cookies
	h.clearSessionDataFromCookies(ctx)

	ctx.JSON(http.StatusOK, resp)
}

// DeleteAllSessions Deletes all sessions of the user
// @Summary Deletes all sessions of the user
// @Description Deletes all sessions of the user
// @Tags    sessions
// @Produce json
// @Id      deleteAllSessions
// @Success 200 {object} sessiondto.DeleteAllSessionsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /sessions/delete-all [post]
func (h *SessionHandler) DeleteAllSessions(ctx *gin.Context) {
	user, err := h.authHandler.GetCurrentUser(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sessiondto.NewDeleteAllSessionsReq()
	req.User = user
	if err = h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sessionUC.DeleteAllSessions(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Clear session cookies
	h.clearSessionDataFromCookies(ctx)

	ctx.JSON(http.StatusOK, resp)
}

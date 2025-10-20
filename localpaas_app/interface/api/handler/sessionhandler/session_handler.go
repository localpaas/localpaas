package sessionhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

type SessionHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	sessionUC   *sessionuc.SessionUC
}

func NewSessionHandler(authHandler *authhandler.AuthHandler, sessionUC *sessionuc.SessionUC) *SessionHandler {
	hdl := &SessionHandler{
		authHandler: authHandler,
		sessionUC:   sessionUC,
	}
	return hdl
}

// GetMe Gets session info of the current user
// @Summary Gets session info of the current user
// @Description Gets session info of the current user
// @Tags    sessions
// @Produce json
// @Id      getMe
// @Success 200 {object} sessiondto.GetMeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /sessions/me [get]
func (h *SessionHandler) GetMe(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sessiondto.NewGetMeReq()
	resp, err := h.sessionUC.GetMe(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RefreshSession Refreshes the current user session
// @Summary Refreshes the current user session
// @Description Refreshes the current user session. Refresh token is required via either
// @Description `Authorization` header or `refresh_token` cookie.
// @Tags    sessions
// @Produce json
// @Id      refreshSession
// @Success 200 {object} sessiondto.RefreshSessionResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /sessions/refresh [post]
func (h *SessionHandler) RefreshSession(ctx *gin.Context) {
	var user *basedto.User
	var err error
	// Refresh token is retrieved from either `Authorization` header or `refresh_token` cookie
	if refreshToken, _ := ctx.Cookie(cookieRefreshToken); refreshToken != "" {
		user, err = h.authHandler.GetCurrentUserByToken(ctx, refreshToken)
	} else {
		user, err = h.authHandler.GetCurrentUser(ctx)
	}
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sessionUC.RefreshSession(h.RequestCtx(ctx), user)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Writes a portion of data response to cookies
	h.writeSessionDataToCookies(ctx, resp.Data.BaseCreateSessionResp, true)

	ctx.JSON(http.StatusOK, resp)
}

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

	resp, err := h.sessionUC.DeleteSession(h.RequestCtx(ctx), user)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Clear session cookies
	h.clearSessionDataFromCookies(ctx)

	ctx.JSON(http.StatusOK, resp)
}

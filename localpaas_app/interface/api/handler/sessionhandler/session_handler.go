package sessionhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/oauthuc"
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

func NewSessionHandler(
	authHandler *authhandler.AuthHandler,
	sessionUC *sessionuc.SessionUC,
	oauthUC *oauthuc.OAuthUC,
) *SessionHandler {
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
// @Param   getAccesses query string false "`getAccesses=true/false`"
// @Success 200 {object} sessiondto.GetMeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /sessions/me [get]
func (h *SessionHandler) GetMe(ctx *gin.Context) {
	user, err := h.authHandler.GetCurrentUser(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sessiondto.NewGetMeReq()
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sessionUC.GetMe(h.RequestCtx(ctx), user, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

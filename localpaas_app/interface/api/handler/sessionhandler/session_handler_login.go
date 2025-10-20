package sessionhandler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

const (
	cookieAccessToken     = "access_token"
	cookieRefreshToken    = "refresh_token"
	cookieAccessHTTPOnly  = false
	cookieRefreshHTTPOnly = true
	cookieRefreshPath     = "/sessions/refresh" // to avoid unnecessary sending of refresh token to server
)

// LoginWithPassword Login to system with username/password
// @Summary Login to system with username/password
// @Description When you get response's `next_step` with value `NextMfa`, you need to call the API
// @Description `/login-with-passcode` to complete the login process.
// @Tags    sessions_auth
// @Produce json
// @Id      loginWithPassword
// @Param   body body sessiondto.LoginWithPasswordReq true "request data"
// @Success 200 {object} sessiondto.LoginWithPasswordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /auth/login-with-password [post]
func (h *SessionHandler) LoginWithPassword(ctx *gin.Context) {
	req := sessiondto.NewLoginWithPasswordReq()
	req.AcceptLanguage = h.ParseRequestLang(ctx)
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sessionUC.LoginWithPassword(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Writes a portion of data response to cookies
	if resp.Data.Session != nil {
		h.writeSessionDataToCookies(ctx, resp.Data.Session, true)
	}

	ctx.JSON(http.StatusOK, resp)
}

// LoginWithPasscode Login to system with passcode after using password
// @Summary Login to system with passcode after using password
// @Description Login to system with passcode after using password
// @Tags    sessions_auth
// @Produce json
// @Id      loginWithPasscode
// @Param   body body sessiondto.LoginWithPasscodeReq true "request data"
// @Success 200 {object} sessiondto.LoginWithPasscodeResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /auth/login-with-passcode [post]
func (h *SessionHandler) LoginWithPasscode(ctx *gin.Context) {
	req := sessiondto.NewLoginWithPasscodeReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sessionUC.LoginWithPasscode(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Writes a portion of data response to cookies
	if resp.Data.Session != nil {
		h.writeSessionDataToCookies(ctx, resp.Data.Session, true)
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *SessionHandler) writeSessionDataToCookies(ctx *gin.Context, sessionResp *sessiondto.BaseCreateSessionResp,
	writeRefreshOnly bool) {
	ctx.SetSameSite(http.SameSiteLaxMode)

	secure := config.Current.IsProdEnv()
	timeNow := timeutil.NowUTC()

	if !writeRefreshOnly {
		accessAge := int(sessionResp.AccessTokenExp.Sub(timeNow).Seconds())
		baseURL := strings.SplitN(config.Current.HTTPServer.BaseURL, ".", 2) //nolint:mnd
		ctx.SetCookie(cookieAccessToken, sessionResp.AccessToken, accessAge, "",
			baseURL[len(baseURL)-1], secure, cookieAccessHTTPOnly)
	}

	refreshPath := ""
	if cookieRefreshPath != "" {
		refreshPath = gofn.Must(url.JoinPath(config.Current.HTTPServer.BasePath, cookieRefreshPath))
	}

	// Writes refresh token only (requested by FE team)
	refreshAge := int(sessionResp.RefreshTokenExp.Sub(timeNow).Seconds())
	ctx.SetCookie(cookieRefreshToken, sessionResp.RefreshToken, refreshAge, refreshPath, "",
		secure, cookieRefreshHTTPOnly)

	// Unsets refresh token fields from the data (maybe unnecessary, but better do it)
	sessionResp.RefreshToken = ""
	sessionResp.RefreshTokenExp = time.Time{}
}

func (h *SessionHandler) clearSessionDataFromCookies(ctx *gin.Context) {
	// MaxAge<0 means delete cookie now
	ctx.SetCookie(cookieAccessToken, "", -1, "", "", false, false)
	ctx.SetCookie(cookieRefreshToken, "", -1, "", "", false, false)
}

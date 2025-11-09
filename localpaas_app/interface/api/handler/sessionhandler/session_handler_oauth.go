package sessionhandler

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc/sessiondto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/oauthuc/oauthdto"
)

const (
	ssoRedirectFEPathOnSuccess = "/auth/sso/success"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// SSOOAuthBegin Starts OAuth SSO flow
// @Summary Starts OAuth SSO flow
// @Description Starts OAuth SSO flow
// @Tags    sessions_auth
// @Produce json
// @Id      ssoOAuthBegin
// @Success 302 "on success redirect to provider OAuth URL"
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /auth/sso/{provider} [get]
func (h *SessionHandler) SSOOAuthBegin(ctx *gin.Context) {
	provider, err := h.ParseStringParam(ctx, "provider")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	oauthType := base.OAuthType(provider)
	if !base.IsValidOAuthType(oauthType) {
		h.RenderError(ctx, apperrors.NewParamInvalid("OAuth").
			WithMsgLog("%s OAuth is not valid", provider))
		return
	}

	oauth, err := h.oauthUC.GetOAuthNoAuth(ctx, &oauthdto.GetOAuthNoAuthReq{
		Name:   provider,
		Status: []base.SettingStatus{base.SettingStatusActive},
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	if oauth == nil {
		h.RenderError(ctx, apperrors.New(apperrors.ErrUnavailable).
			WithMsgLog("%s OAuth SSO is not configured", provider))
		return
	}

	baseCallbackURL := h.oauthUC.GetOAuthBaseCallbackURL()
	switch oauthType {
	case base.OAuthTypeGithub:
		goth.UseProviders(github.New(oauth.ClientID, oauth.ClientSecret, baseCallbackURL+"/"+provider, oauth.Scopes...))
	case base.OAuthTypeGitlab:
		goth.UseProviders(gitlab.New(oauth.ClientID, oauth.ClientSecret, baseCallbackURL+"/"+provider, oauth.Scopes...))
	case base.OAuthTypeGoogle:
		goth.UseProviders(google.New(oauth.ClientID, oauth.ClientSecret, baseCallbackURL+"/"+provider, oauth.Scopes...))

	// Custom types
	case base.OAuthTypeGitlabCustom:
		providerObj := gitlab.NewCustomisedURL(oauth.ClientID, oauth.ClientSecret,
			baseCallbackURL+"/"+provider, oauth.AuthURL, oauth.TokenURL, oauth.ProfileURL,
			oauth.Scopes...)
		providerObj.SetName(provider)
		goth.UseProviders(providerObj)
	}

	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// SSOOAuthCallback Begins SSO flow
// @Summary Begins SSO flow
// @Description Begins SSO flow
// @Tags    users
// @Produce json
// @Id      beginOAuth
// @Param   provider path string true "provider name"
// @Success 302 "on success redirect to the dashboard page"
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /auth/sso/callback/{provider} [get]
// @Router  /auth/sso/callback/{provider} [post]
func (h *SessionHandler) SSOOAuthCallback(ctx *gin.Context) {
	provider, err := h.ParseStringParam(ctx, "provider")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	oauthUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Create a new session for OAuth user in our service
	sessionReq := sessiondto.NewCreateOAuthSessionReq()
	sessionReq.User = &oauthUser
	sessionResp, err := h.sessionUC.CreateOAuthSession(h.RequestCtx(ctx), sessionReq)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Write session data to request cookies as we will redirect to a front-end path
	// and passing long tokens in the URL is risky.
	h.writeSessionDataToCookies(ctx, &sessionResp.BaseCreateSessionResp, false)

	// Redirect client to front-end page
	redirectURL := gofn.Must(url.JoinPath(config.Current.App.BaseURL, ssoRedirectFEPathOnSuccess))
	ctx.Redirect(http.StatusFound, redirectURL)
}

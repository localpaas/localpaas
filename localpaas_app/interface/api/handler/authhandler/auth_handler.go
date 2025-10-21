package authhandler

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sessionuc"
)

var (
	NoAccessCheck = (*permission.AccessCheck)(nil)
)

type AuthHandler struct {
	*handler.BaseHandler
	sessionUC *sessionuc.SessionUC
}

func NewAuthHandler(
	sessionUC *sessionuc.SessionUC,
) *AuthHandler {
	hdl := &AuthHandler{
		sessionUC: sessionUC,
	}
	return hdl
}

func (h *AuthHandler) GetCurrentUser(ctx *gin.Context) (*basedto.User, error) {
	token, err := h.getAuthToken(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}

	user, err := h.sessionUC.GetCurrentUser(h.RequestCtx(ctx), token)
	if err != nil {
		return nil, apperrors.New(err)
	}

	return user, nil
}

func (h *AuthHandler) GetCurrentUserByToken(ctx *gin.Context, token string) (*basedto.User, error) {
	user, err := h.sessionUC.GetCurrentUser(h.RequestCtx(ctx), token)
	if err != nil {
		return nil, apperrors.New(err)
	}

	return user, nil
}

func (h *AuthHandler) GetCurrentAuth(ctx *gin.Context, accessCheck *permission.AccessCheck) (*basedto.Auth, error) {
	auth, err := h.getCurrentAuth(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}

	//nolint:staticcheck
	if accessCheck != NoAccessCheck && !(config.Current.IsDevEnv() && config.Current.DevMode.SkipAuthCheck) {
		if err = h.sessionUC.VerifyAuth(ctx, auth, accessCheck); err != nil {
			// NOTE: even on error, we still return the `auth` object so the client code
			// still can be able to check permission with another method.
			return auth, apperrors.New(err)
		}
	}

	return auth, nil
}

func (h *AuthHandler) getCurrentAuth(ctx *gin.Context) (*basedto.Auth, error) {
	token, err := h.getAuthToken(ctx)
	if err != nil {
		return nil, apperrors.New(err)
	}

	auth, err := h.sessionUC.GetCurrentAuth(h.RequestCtx(ctx), token)
	if err != nil {
		return nil, apperrors.New(err)
	}
	return auth, nil
}

// getAuthToken gets token from request header `Authorization`.
// The value should be in form of `Bearer <token-data>`.
func (h *AuthHandler) getAuthToken(ctx *gin.Context) (token string, err error) {
	tokenParts := strings.SplitN(ctx.GetHeader("Authorization"), " ", 2) //nolint:mnd
	if len(tokenParts) != 2 || tokenParts[1] == "" {
		return "", apperrors.New(apperrors.ErrSessionJWTInvalid)
	}
	return tokenParts[1], nil
}

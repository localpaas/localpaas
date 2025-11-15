package usersettingshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/usersettings/apikeyuc/apikeydto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListAPIKey Lists API key
// @Summary Lists API key
// @Description Lists API key
// @Tags    user_settings_api_keys
// @Produce json
// @Id      listAPIKeySettings
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} apikeydto.ListAPIKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/current/settings/api-keys [get]
func (h *UserSettingsHandler) ListAPIKey(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := apikeydto.NewListAPIKeyReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.apiKeyUC.ListAPIKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetAPIKey Gets API key details
// @Summary Gets API key details
// @Description Gets API key details
// @Tags    user_settings_api_keys
// @Produce json
// @Id      getAPIKeySetting
// @Param   ID path string true "s3 storage ID"
// @Success 200 {object} apikeydto.GetAPIKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/current/settings/api-keys/{ID} [get]
func (h *UserSettingsHandler) GetAPIKey(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := apikeydto.NewGetAPIKeyReq()
	req.ID = id
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.apiKeyUC.GetAPIKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateAPIKey Creates a new API key
// @Summary Creates a new API key
// @Description Creates a new API key
// @Tags    user_settings_api_keys
// @Produce json
// @Id      createAPIKeySetting
// @Param   body body apikeydto.CreateAPIKeyReq true "request data"
// @Success 201 {object} apikeydto.CreateAPIKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/current/settings/api-keys [post]
func (h *UserSettingsHandler) CreateAPIKey(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Not allow to use API key to create API key
	if auth.User.AuthClaims.IsAPIKey {
		h.RenderError(ctx, apperrors.New(apperrors.ErrForbidden).
			WithMsgLog("not allow to create API key by using API key session"))
		return
	}

	req := apikeydto.NewCreateAPIKeyReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.apiKeyUC.CreateAPIKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteAPIKey Deletes an API key
// @Summary Deletes an API key
// @Description Deletes an API key
// @Tags    user_settings_api_keys
// @Produce json
// @Id      deleteAPIKeySetting
// @Param   ID path string true "API key ID"
// @Success 200 {object} apikeydto.DeleteAPIKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /users/current/settings/api-keys/{ID} [delete]
func (h *UserSettingsHandler) DeleteAPIKey(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Not allow to use API key to delete API key
	if auth.User.AuthClaims.IsAPIKey {
		h.RenderError(ctx, apperrors.New(apperrors.ErrForbidden).
			WithMsgLog("not allow to delete API key by using API key session"))
		return
	}

	req := apikeydto.NewDeleteAPIKeyReq()
	req.ID = id
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.apiKeyUC.DeleteAPIKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

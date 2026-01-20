package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc/discorddto"
)

// ListDiscord Lists Discord providers
// @Summary Lists Discord providers
// @Description Lists Discord providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderDiscord
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} discorddto.ListDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord [get]
func (h *ProvidersHandler) ListDiscord(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeGlobal)
}

// GetDiscord Gets Discord provider details
// @Summary Gets Discord provider details
// @Description Gets Discord provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderDiscord
// @Param   id path string true "provider ID"
// @Success 200 {object} discorddto.GetDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{id} [get]
func (h *ProvidersHandler) GetDiscord(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeGlobal)
}

// CreateDiscord Creates a new Discord provider
// @Summary Creates a new Discord provider
// @Description Creates a new Discord provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderDiscord
// @Param   body body discorddto.CreateDiscordReq true "request data"
// @Success 201 {object} discorddto.CreateDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord [post]
func (h *ProvidersHandler) CreateDiscord(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeGlobal)
}

// UpdateDiscord Updates Discord provider
// @Summary Updates Discord provider
// @Description Updates Discord provider
// @Tags    global_providers
// @Produce json
// @Id      updateProviderDiscord
// @Param   id path string true "provider ID"
// @Param   body body discorddto.UpdateDiscordReq true "request data"
// @Success 200 {object} discorddto.UpdateDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{id} [put]
func (h *ProvidersHandler) UpdateDiscord(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeGlobal)
}

// UpdateDiscordMeta Updates Discord meta provider
// @Summary Updates Discord meta provider
// @Description Updates Discord meta provider
// @Tags    global_providers
// @Produce json
// @Id      updateProviderDiscordMeta
// @Param   id path string true "provider ID"
// @Param   body body discorddto.UpdateDiscordMetaReq true "request data"
// @Success 200 {object} discorddto.UpdateDiscordMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{id}/meta [put]
func (h *ProvidersHandler) UpdateDiscordMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeDiscord, base.SettingScopeGlobal)
}

// DeleteDiscord Deletes Discord provider
// @Summary Deletes Discord provider
// @Description Deletes Discord provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderDiscord
// @Param   id path string true "provider ID"
// @Success 200 {object} discorddto.DeleteDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{id} [delete]
func (h *ProvidersHandler) DeleteDiscord(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeGlobal)
}

// TestSendDiscordMsg Tests sending a msg
// @Summary Tests sending a msg
// @Description Tests sending a msg
// @Tags    global_providers
// @Produce json
// @Id      testSendDiscordMsg
// @Param   body body discorddto.TestSendDiscordMsgReq true "request data"
// @Success 200 {object} discorddto.TestSendDiscordMsgResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/test-send-msg [post]
func (h *ProvidersHandler) TestSendDiscordMsg(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewTestSendDiscordMsgReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.DiscordUC.TestSendDiscordMsg(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

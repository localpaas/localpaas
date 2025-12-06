package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListDiscord Lists Discord providers
// @Summary Lists Discord providers
// @Description Lists Discord providers
// @Tags    providers_discord
// @Produce json
// @Id      listDiscordProviders
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} discorddto.ListDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord [get]
func (h *ProvidersHandler) ListDiscord(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeDiscord,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewListDiscordReq()
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.discordUC.ListDiscord(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetDiscord Gets Discord provider details
// @Summary Gets Discord provider details
// @Description Gets Discord provider details
// @Tags    providers_discord
// @Produce json
// @Id      getDiscordProvider
// @Param   ID path string true "provider ID"
// @Success 200 {object} discorddto.GetDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{ID} [get]
func (h *ProvidersHandler) GetDiscord(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeDiscord,
		ResourceID:     id,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewGetDiscordReq()
	req.ID = id
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.discordUC.GetDiscord(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateDiscord Creates a new Discord provider
// @Summary Creates a new Discord provider
// @Description Creates a new Discord provider
// @Tags    providers_discord
// @Produce json
// @Id      createDiscordProvider
// @Param   body body discorddto.CreateDiscordReq true "request data"
// @Success 201 {object} discorddto.CreateDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord [post]
func (h *ProvidersHandler) CreateDiscord(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewCreateDiscordReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.discordUC.CreateDiscord(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateDiscord Updates Discord provider
// @Summary Updates Discord provider
// @Description Updates Discord provider
// @Tags    providers_discord
// @Produce json
// @Id      updateDiscordProvider
// @Param   ID path string true "provider ID"
// @Param   body body discorddto.UpdateDiscordReq true "request data"
// @Success 200 {object} discorddto.UpdateDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{ID} [put]
func (h *ProvidersHandler) UpdateDiscord(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeDiscord,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewUpdateDiscordReq()
	req.ID = id
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.discordUC.UpdateDiscord(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateDiscordMeta Updates Discord meta provider
// @Summary Updates Discord meta provider
// @Description Updates Discord meta provider
// @Tags    providers_discord
// @Produce json
// @Id      updateDiscordProviderMeta
// @Param   ID path string true "provider ID"
// @Param   body body discorddto.UpdateDiscordMetaReq true "request data"
// @Success 200 {object} discorddto.UpdateDiscordMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{ID}/meta [put]
func (h *ProvidersHandler) UpdateDiscordMeta(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeDiscord,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewUpdateDiscordMetaReq()
	req.ID = id
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.discordUC.UpdateDiscordMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteDiscord Deletes Discord provider
// @Summary Deletes Discord provider
// @Description Deletes Discord provider
// @Tags    providers_discord
// @Produce json
// @Id      deleteDiscordProvider
// @Param   ID path string true "provider ID"
// @Success 200 {object} discorddto.DeleteDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/discord/{ID} [delete]
func (h *ProvidersHandler) DeleteDiscord(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeDiscord,
		ResourceID:     id,
		Action:         base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewDeleteDiscordReq()
	req.ID = id
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.discordUC.DeleteDiscord(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// TestSendDiscordMsg Tests sending a msg
// @Summary Tests sending a msg
// @Description Tests sending a msg
// @Tags    providers_discord
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

	resp, err := h.discordUC.TestSendDiscordMsg(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

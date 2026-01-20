package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/discorduc/discorddto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListDiscord Lists Discord providers
// @Summary Lists Discord providers
// @Description Lists Discord providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectDiscord
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} discorddto.ListDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/discord [get]
func (h *ProjectHandler) ListDiscord(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewListDiscordReq()
	req.ProjectID = projectID
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
// @Tags    project_providers
// @Produce json
// @Id      getProjectDiscord
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} discorddto.GetDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/discord/{id} [get]
func (h *ProjectHandler) GetDiscord(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewGetDiscordReq()
	req.ID = id
	req.ProjectID = projectID
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
// @Tags    project_providers
// @Produce json
// @Id      createProjectDiscord
// @Param   projectID path string true "project ID"
// @Param   body body discorddto.CreateDiscordReq true "request data"
// @Success 201 {object} discorddto.CreateDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/discord [post]
func (h *ProjectHandler) CreateDiscord(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewCreateDiscordReq()
	req.ProjectID = projectID
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
// @Tags    project_providers
// @Produce json
// @Id      updateProjectDiscord
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body discorddto.UpdateDiscordReq true "request data"
// @Success 200 {object} discorddto.UpdateDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/discord/{id} [put]
func (h *ProjectHandler) UpdateDiscord(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewUpdateDiscordReq()
	req.ID = id
	req.ProjectID = projectID
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
// @Tags    project_providers
// @Produce json
// @Id      updateProjectDiscordMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body discorddto.UpdateDiscordMetaReq true "request data"
// @Success 200 {object} discorddto.UpdateDiscordMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/discord/{id}/meta [put]
func (h *ProjectHandler) UpdateDiscordMeta(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewUpdateDiscordMetaReq()
	req.ID = id
	req.ProjectID = projectID
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
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectDiscord
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} discorddto.DeleteDiscordResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/discord/{id} [delete]
func (h *ProjectHandler) DeleteDiscord(ctx *gin.Context) {
	auth, projectID, id, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := discorddto.NewDeleteDiscordReq()
	req.ID = id
	req.ProjectID = projectID
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

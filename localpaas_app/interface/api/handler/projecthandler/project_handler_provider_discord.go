package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/discorduc/discorddto"
)

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
	h.ListSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeProject)
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
	h.GetSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeProject)
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
	h.CreateSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeProject)
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
	h.UpdateSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeProject)
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
	h.UpdateSettingMeta(ctx, base.ResourceTypeDiscord, base.SettingScopeProject)
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
	h.DeleteSetting(ctx, base.ResourceTypeDiscord, base.SettingScopeProject)
}

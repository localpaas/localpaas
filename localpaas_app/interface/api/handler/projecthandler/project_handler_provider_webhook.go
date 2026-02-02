package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/webhookuc/webhookdto"
)

// ListWebhook Lists webhook providers
// @Summary Lists webhook providers
// @Description Lists webhook providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectWebhook
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} webhookdto.ListWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/webhooks [get]
func (h *ProjectHandler) ListWebhook(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeProject)
}

// GetWebhook Gets webhook provider details
// @Summary Gets webhook provider details
// @Description Gets webhook provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectWebhook
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} webhookdto.GetWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/webhooks/{id} [get]
func (h *ProjectHandler) GetWebhook(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeProject)
}

// CreateWebhook Creates a new webhook provider
// @Summary Creates a new webhook provider
// @Description Creates a new webhook provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectWebhook
// @Param   projectID path string true "project ID"
// @Param   body body webhookdto.CreateWebhookReq true "request data"
// @Success 201 {object} webhookdto.CreateWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/webhooks [post]
func (h *ProjectHandler) CreateWebhook(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeProject)
}

// UpdateWebhook Updates webhook
// @Summary Updates webhook
// @Description Updates webhook
// @Tags    project_providers
// @Produce json
// @Id      updateProjectWebhook
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body webhookdto.UpdateWebhookReq true "request data"
// @Success 200 {object} webhookdto.UpdateWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/webhooks/{id} [put]
func (h *ProjectHandler) UpdateWebhook(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeProject)
}

// UpdateWebhookMeta Updates webhook meta
// @Summary Updates webhook meta
// @Description Updates webhook meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectWebhookMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body webhookdto.UpdateWebhookMetaReq true "request data"
// @Success 200 {object} webhookdto.UpdateWebhookMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/webhooks/{id}/meta [put]
func (h *ProjectHandler) UpdateWebhookMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeWebhook, base.SettingScopeProject)
}

// DeleteWebhook Deletes webhook provider
// @Summary Deletes webhook provider
// @Description Deletes webhook provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectWebhook
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} webhookdto.DeleteWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/webhooks/{id} [delete]
func (h *ProjectHandler) DeleteWebhook(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeProject)
}

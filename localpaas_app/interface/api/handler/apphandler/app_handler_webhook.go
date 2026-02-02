package apphandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/webhookuc/webhookdto"
)

// ListAppWebhook Lists webhooks
// @Summary Lists webhooks
// @Description Lists webhooks
// @Tags    apps
// @Produce json
// @Id      listAppWebhook
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} webhookdto.ListWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/webhooks [get]
func (h *AppHandler) ListAppWebhook(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeApp)
}

// GetAppWebhook Gets webhook details
// @Summary Gets webhook details
// @Description Gets webhook details
// @Tags    apps
// @Produce json
// @Id      getAppWebhook
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} webhookdto.GetWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/webhooks/{id} [get]
func (h *AppHandler) GetAppWebhook(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeApp)
}

// CreateAppWebhook Creates a new webhook
// @Summary Creates a new webhook
// @Description Creates a new webhook
// @Tags    apps
// @Produce json
// @Id      createAppWebhook
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body webhookdto.CreateWebhookReq true "request data"
// @Success 201 {object} webhookdto.CreateWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/webhooks [post]
func (h *AppHandler) CreateAppWebhook(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeApp)
}

// UpdateAppWebhook Updates a webhook
// @Summary Updates a webhook
// @Description Updates a webhook
// @Tags    apps
// @Produce json
// @Id      updateAppWebhook
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Param   body body webhookdto.UpdateWebhookReq true "request data"
// @Success 200 {object} webhookdto.UpdateWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/webhooks/{id} [put]
func (h *AppHandler) UpdateAppWebhook(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeApp)
}

// UpdateAppWebhookMeta Updates webhook meta
// @Summary Updates webhook meta
// @Description Updates webhook meta
// @Tags    apps
// @Produce json
// @Id      updateAppWebhookMeta
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Param   body body webhookdto.UpdateWebhookMetaReq true "request data"
// @Success 200 {object} webhookdto.UpdateWebhookMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/webhooks/{id}/meta [put]
func (h *AppHandler) UpdateAppWebhookMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeWebhook, base.SettingScopeApp)
}

// DeleteAppWebhook Deletes webhook
// @Summary Deletes webhook
// @Description Deletes webhook
// @Tags    apps
// @Produce json
// @Id      deleteAppWebhook
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} webhookdto.DeleteWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/webhooks/{id} [delete]
func (h *AppHandler) DeleteAppWebhook(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeWebhook, base.SettingScopeApp)
}

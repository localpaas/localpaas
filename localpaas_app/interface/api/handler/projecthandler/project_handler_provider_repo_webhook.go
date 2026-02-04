package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/repowebhookuc/repowebhookdto"
)

// ListRepoWebhook Lists webhook providers
// @Summary Lists webhook providers
// @Description Lists webhook providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectRepoWebhook
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} repowebhookdto.ListRepoWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/repo-webhooks [get]
func (h *ProjectHandler) ListRepoWebhook(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeRepoWebhook, base.SettingScopeProject)
}

// GetRepoWebhook Gets webhook provider details
// @Summary Gets webhook provider details
// @Description Gets webhook provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectRepoWebhook
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} repowebhookdto.GetRepoWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/repo-webhooks/{id} [get]
func (h *ProjectHandler) GetRepoWebhook(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeRepoWebhook, base.SettingScopeProject)
}

// CreateRepoWebhook Creates a new webhook provider
// @Summary Creates a new webhook provider
// @Description Creates a new webhook provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectRepoWebhook
// @Param   projectID path string true "project ID"
// @Param   body body repowebhookdto.CreateRepoWebhookReq true "request data"
// @Success 201 {object} repowebhookdto.CreateRepoWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/repo-webhooks [post]
func (h *ProjectHandler) CreateRepoWebhook(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeRepoWebhook, base.SettingScopeProject)
}

// UpdateRepoWebhook Updates webhook
// @Summary Updates webhook
// @Description Updates webhook
// @Tags    project_providers
// @Produce json
// @Id      updateProjectRepoWebhook
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body repowebhookdto.UpdateRepoWebhookReq true "request data"
// @Success 200 {object} repowebhookdto.UpdateRepoWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/repo-webhooks/{id} [put]
func (h *ProjectHandler) UpdateRepoWebhook(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeRepoWebhook, base.SettingScopeProject)
}

// UpdateRepoWebhookMeta Updates webhook meta
// @Summary Updates webhook meta
// @Description Updates webhook meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectRepoWebhookMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body repowebhookdto.UpdateRepoWebhookMetaReq true "request data"
// @Success 200 {object} repowebhookdto.UpdateRepoWebhookMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/repo-webhooks/{id}/meta [put]
func (h *ProjectHandler) UpdateRepoWebhookMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeRepoWebhook, base.SettingScopeProject)
}

// DeleteRepoWebhook Deletes webhook provider
// @Summary Deletes webhook provider
// @Description Deletes webhook provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectRepoWebhook
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} repowebhookdto.DeleteRepoWebhookResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/repo-webhooks/{id} [delete]
func (h *ProjectHandler) DeleteRepoWebhook(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeRepoWebhook, base.SettingScopeProject)
}

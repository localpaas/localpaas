package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// ListSsl Lists SSL providers
// @Summary Lists SSL providers
// @Description Lists SSL providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectSSL
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls [get]
func (h *ProjectHandler) ListSsl(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// GetSsl Gets SSL provider details
// @Summary Gets SSL provider details
// @Description Gets SSL provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectSSL
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} ssldto.GetSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id} [get]
func (h *ProjectHandler) GetSsl(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// CreateSsl Creates a new SSL provider
// @Summary Creates a new SSL provider
// @Description Creates a new SSL provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectSSL
// @Param   projectID path string true "project ID"
// @Param   body body ssldto.CreateSslReq true "request data"
// @Success 201 {object} ssldto.CreateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls [post]
func (h *ProjectHandler) CreateSsl(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// UpdateSsl Updates SSL
// @Summary Updates SSL
// @Description Updates SSL
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSSL
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body ssldto.UpdateSslReq true "request data"
// @Success 200 {object} ssldto.UpdateSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id} [put]
func (h *ProjectHandler) UpdateSsl(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// UpdateSslMeta Updates SSL meta
// @Summary Updates SSL meta
// @Description Updates SSL meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSSLMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body ssldto.UpdateSslMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSslMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id}/meta [put]
func (h *ProjectHandler) UpdateSslMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// DeleteSsl Deletes SSL provider
// @Summary Deletes SSL provider
// @Description Deletes SSL provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectSSL
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} ssldto.DeleteSslResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssls/{id} [delete]
func (h *ProjectHandler) DeleteSsl(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

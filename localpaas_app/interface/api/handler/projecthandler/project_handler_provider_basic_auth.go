package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/basicauthuc/basicauthdto"
)

// ListBasicAuth Lists basic auth providers
// @Summary Lists basic auth providers
// @Description Lists basic auth providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} basicauthdto.ListBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth [get]
func (h *ProjectHandler) ListBasicAuth(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeProject)
}

// GetBasicAuth Gets basic auth provider details
// @Summary Gets basic auth provider details
// @Description Gets basic auth provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} basicauthdto.GetBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{itemID} [get]
func (h *ProjectHandler) GetBasicAuth(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeProject)
}

// CreateBasicAuth Creates a new basic auth provider
// @Summary Creates a new basic auth provider
// @Description Creates a new basic auth provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   body body basicauthdto.CreateBasicAuthReq true "request data"
// @Success 201 {object} basicauthdto.CreateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth [post]
func (h *ProjectHandler) CreateBasicAuth(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeProject)
}

// UpdateBasicAuth Updates basic auth
// @Summary Updates basic auth
// @Description Updates basic auth
// @Tags    project_providers
// @Produce json
// @Id      updateProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body basicauthdto.UpdateBasicAuthReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{itemID} [put]
func (h *ProjectHandler) UpdateBasicAuth(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeProject)
}

// UpdateBasicAuthMeta Updates basic auth meta
// @Summary Updates basic auth meta
// @Description Updates basic auth meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectBasicAuthMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body basicauthdto.UpdateBasicAuthMetaReq true "request data"
// @Success 200 {object} basicauthdto.UpdateBasicAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{itemID}/meta [put]
func (h *ProjectHandler) UpdateBasicAuthMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeBasicAuth, base.SettingScopeProject)
}

// DeleteBasicAuth Deletes basic auth provider
// @Summary Deletes basic auth provider
// @Description Deletes basic auth provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectBasicAuth
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} basicauthdto.DeleteBasicAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/basic-auth/{itemID} [delete]
func (h *ProjectHandler) DeleteBasicAuth(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeBasicAuth, base.SettingScopeProject)
}

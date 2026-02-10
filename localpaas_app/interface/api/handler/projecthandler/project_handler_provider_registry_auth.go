package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

// ListRegistryAuth Lists registry auth providers
// @Summary Lists registry auth providers
// @Description Lists registry auth providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} registryauthdto.ListRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth [get]
func (h *ProjectHandler) ListRegistryAuth(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeProject)
}

// GetRegistryAuth Gets registry auth provider details
// @Summary Gets registry auth provider details
// @Description Gets registry auth provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} registryauthdto.GetRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{itemID} [get]
func (h *ProjectHandler) GetRegistryAuth(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeProject)
}

// CreateRegistryAuth Creates a new registry auth provider
// @Summary Creates a new registry auth provider
// @Description Creates a new registry auth provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   body body registryauthdto.CreateRegistryAuthReq true "request data"
// @Success 201 {object} registryauthdto.CreateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth [post]
func (h *ProjectHandler) CreateRegistryAuth(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeProject)
}

// UpdateRegistryAuth Updates registry auth
// @Summary Updates registry auth
// @Description Updates registry auth
// @Tags    project_providers
// @Produce json
// @Id      updateProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body registryauthdto.UpdateRegistryAuthReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{itemID} [put]
func (h *ProjectHandler) UpdateRegistryAuth(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeProject)
}

// UpdateRegistryAuthMeta Updates registry auth meta
// @Summary Updates registry auth meta
// @Description Updates registry auth meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectRegistryAuthMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body registryauthdto.UpdateRegistryAuthMetaReq true "request data"
// @Success 200 {object} registryauthdto.UpdateRegistryAuthMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{itemID}/meta [put]
func (h *ProjectHandler) UpdateRegistryAuthMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeProject)
}

// DeleteRegistryAuth Deletes registry auth provider
// @Summary Deletes registry auth provider
// @Description Deletes registry auth provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectRegistryAuth
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} registryauthdto.DeleteRegistryAuthResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/registry-auth/{itemID} [delete]
func (h *ProjectHandler) DeleteRegistryAuth(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeRegistryAuth, base.SettingScopeProject)
}

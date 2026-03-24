package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
)

// ListCloudProvider Lists cloud providers
// @Summary Lists cloud providers
// @Description Lists cloud providers
// @Tags    project_settings
// @Produce json
// @Id      listProjectCloudProvider
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} cloudproviderdto.ListCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cloud-providers [get]
func (h *ProjectHandler) ListCloudProvider(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeProject)
}

// GetCloudProvider Gets cloud provider details
// @Summary Gets cloud provider details
// @Description Gets cloud provider details
// @Tags    project_settings
// @Produce json
// @Id      getProjectCloudProvider
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cloudproviderdto.GetCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cloud-providers/{itemID} [get]
func (h *ProjectHandler) GetCloudProvider(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeProject)
}

// CreateCloudProvider Creates a cloud provider
// @Summary Creates a cloud provider
// @Description Creates a cloud provider
// @Tags    project_settings
// @Produce json
// @Id      createProjectCloudProvider
// @Param   projectID path string true "project ID"
// @Param   body body cloudproviderdto.CreateCloudProviderReq true "request data"
// @Success 201 {object} cloudproviderdto.CreateCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cloud-providers [post]
func (h *ProjectHandler) CreateCloudProvider(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeProject)
}

// UpdateCloudProvider Updates a cloud provider
// @Summary Updates a cloud provider
// @Description Updates a cloud provider
// @Tags    project_settings
// @Produce json
// @Id      updateProjectCloudProvider
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body cloudproviderdto.UpdateCloudProviderReq true "request data"
// @Success 200 {object} cloudproviderdto.UpdateCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cloud-providers/{itemID} [put]
func (h *ProjectHandler) UpdateCloudProvider(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeProject)
}

// UpdateCloudProviderMeta Updates a cloud provider's meta
// @Summary Updates a cloud provider's meta
// @Description Updates a cloud provider's meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectCloudProviderMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body cloudproviderdto.UpdateCloudProviderMetaReq true "request data"
// @Success 200 {object} cloudproviderdto.UpdateCloudProviderMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cloud-providers/{itemID}/meta [put]
func (h *ProjectHandler) UpdateCloudProviderMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeCloudProvider, base.SettingScopeProject)
}

// DeleteCloudProvider Deletes a cloud provider
// @Summary Deletes a cloud provider
// @Description Deletes a cloud provider
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectCloudProvider
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cloudproviderdto.DeleteCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/cloud-providers/{itemID} [delete]
func (h *ProjectHandler) DeleteCloudProvider(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeProject)
}

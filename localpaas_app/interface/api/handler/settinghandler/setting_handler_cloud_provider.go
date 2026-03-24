package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudprovideruc/cloudproviderdto"
)

// ListCloudProvider Lists cloud providers
// @Summary Lists cloud providers
// @Description Lists cloud providers
// @Tags    settings
// @Produce json
// @Id      listSettingCloudProvider
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} cloudproviderdto.ListCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-providers [get]
func (h *SettingHandler) ListCloudProvider(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeGlobal)
}

// GetCloudProvider Gets cloud provider details
// @Summary Gets cloud provider details
// @Description Gets cloud provider details
// @Tags    settings
// @Produce json
// @Id      getSettingCloudProvider
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cloudproviderdto.GetCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-providers/{itemID} [get]
func (h *SettingHandler) GetCloudProvider(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeGlobal)
}

// CreateCloudProvider Creates a new cloud provider
// @Summary Creates a new cloud provider
// @Description Creates a new cloud provider
// @Tags    settings
// @Produce json
// @Id      createSettingCloudProvider
// @Param   body body cloudproviderdto.CreateCloudProviderReq true "request data"
// @Success 201 {object} cloudproviderdto.CreateCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-providers [post]
func (h *SettingHandler) CreateCloudProvider(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeGlobal)
}

// UpdateCloudProvider Updates a cloud provider
// @Summary Updates a cloud provider
// @Description Updates a cloud provider
// @Tags    settings
// @Produce json
// @Id      updateSettingCloudProvider
// @Param   itemID path string true "setting ID"
// @Param   body body cloudproviderdto.UpdateCloudProviderReq true "request data"
// @Success 200 {object} cloudproviderdto.UpdateCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-providers/{itemID} [put]
func (h *SettingHandler) UpdateCloudProvider(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeGlobal)
}

// UpdateCloudProviderMeta Updates a cloud provider's metadata
// @Summary Updates a cloud provider's metadata
// @Description Updates a cloud provider's metadata
// @Tags    settings
// @Produce json
// @Id      updateSettingCloudProviderMeta
// @Param   itemID path string true "setting ID"
// @Param   body body cloudproviderdto.UpdateCloudProviderMetaReq true "request data"
// @Success 200 {object} cloudproviderdto.UpdateCloudProviderMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-providers/{itemID}/meta [put]
func (h *SettingHandler) UpdateCloudProviderMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeCloudProvider, base.SettingScopeGlobal)
}

// DeleteCloudProvider Deletes a cloud provider
// @Summary Deletes a cloud provider
// @Description Deletes a cloud provider
// @Tags    settings
// @Produce json
// @Id      deleteSettingCloudProvider
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cloudproviderdto.DeleteCloudProviderResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-providers/{itemID} [delete]
func (h *SettingHandler) DeleteCloudProvider(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeCloudProvider, base.SettingScopeGlobal)
}

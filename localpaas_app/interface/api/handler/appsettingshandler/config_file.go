package appsettingshandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/configfileuc/configfiledto"
)

// ListConfigFile Lists app config files
// @Summary Lists app config files
// @Description Lists app config files
// @Tags    apps
// @Produce json
// @Id      listAppConfigFile
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} configfiledto.ListConfigFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/config-files [get]
func (h *Handler) ListConfigFile(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeConfigFile, base.SettingScopeApp)
}

// GetConfigFile Get an app config file details
// @Summary Get an app config file details
// @Description Get an app config file details
// @Tags    apps
// @Produce json
// @Id      getAppConfigFile
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} configfiledto.GetConfigFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/config-files/{itemID} [get]
func (h *Handler) GetConfigFile(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeConfigFile, base.SettingScopeApp)
}

// CreateConfigFile Creates an app config file
// @Summary Creates an app config file
// @Description Creates an app config file
// @Tags    apps
// @Produce json
// @Id      createAppConfigFile
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body configfiledto.CreateConfigFileReq true "request data"
// @Success 201 {object} configfiledto.CreateConfigFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/config-files [post]
func (h *Handler) CreateConfigFile(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeConfigFile, base.SettingScopeApp)
}

// UpdateConfigFile Updates an app config file
// @Summary Updates an app config file
// @Description Updates an app config file
// @Tags    apps
// @Produce json
// @Id      updateAppConfigFile
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Param   body body configfiledto.UpdateConfigFileReq true "request data"
// @Success 200 {object} configfiledto.UpdateConfigFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/config-files/{itemID} [put]
func (h *Handler) UpdateConfigFile(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeConfigFile, base.SettingScopeApp)
}

// UpdateConfigFileStatus Updates app config file status
// @Summary Updates app config file status
// @Description Updates app config file status
// @Tags    apps
// @Produce json
// @Id      updateAppConfigFileStatus
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Param   body body configfiledto.UpdateConfigFileStatusReq true "request data"
// @Success 200 {object} configfiledto.UpdateConfigFileStatusResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/config-files/{itemID}/status [put]
func (h *Handler) UpdateConfigFileStatus(ctx *gin.Context) {
	h.UpdateSettingStatus(ctx, base.ResourceTypeConfigFile, base.SettingScopeApp)
}

// DeleteConfigFile Deletes an app config file
// @Summary Deletes an app config file
// @Description Deletes an app config file
// @Tags    apps
// @Produce json
// @Id      deleteAppConfigFile
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} configfiledto.DeleteConfigFileResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/config-files/{itemID} [delete]
func (h *Handler) DeleteConfigFile(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeConfigFile, base.SettingScopeApp)
}

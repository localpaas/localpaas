package apphandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/healthcheckuc/healthcheckdto"
)

// ListAppHealthcheck Lists healthchecks
// @Summary Lists healthchecks
// @Description Lists healthchecks
// @Tags    apps
// @Produce json
// @Id      listAppHealthcheck
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} healthcheckdto.ListHealthcheckResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/healthchecks [get]
func (h *AppHandler) ListAppHealthcheck(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeHealthcheck, base.SettingScopeApp)
}

// GetAppHealthcheck Gets healthcheck details
// @Summary Gets healthcheck details
// @Description Gets healthcheck details
// @Tags    apps
// @Produce json
// @Id      getAppHealthcheck
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} healthcheckdto.GetHealthcheckResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/healthchecks/{itemID} [get]
func (h *AppHandler) GetAppHealthcheck(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeHealthcheck, base.SettingScopeApp)
}

// CreateAppHealthcheck Creates a new healthcheck
// @Summary Creates a new healthcheck
// @Description Creates a new healthcheck
// @Tags    apps
// @Produce json
// @Id      createAppHealthcheck
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body healthcheckdto.CreateHealthcheckReq true "request data"
// @Success 201 {object} healthcheckdto.CreateHealthcheckResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/healthchecks [post]
func (h *AppHandler) CreateAppHealthcheck(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeHealthcheck, base.SettingScopeApp)
}

// UpdateAppHealthcheck Updates a healthcheck
// @Summary Updates a healthcheck
// @Description Updates a healthcheck
// @Tags    apps
// @Produce json
// @Id      updateAppHealthcheck
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Param   body body healthcheckdto.UpdateHealthcheckReq true "request data"
// @Success 200 {object} healthcheckdto.UpdateHealthcheckResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/healthchecks/{itemID} [put]
func (h *AppHandler) UpdateAppHealthcheck(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeHealthcheck, base.SettingScopeApp)
}

// UpdateAppHealthcheckMeta Updates healthcheck meta
// @Summary Updates healthcheck meta
// @Description Updates healthcheck meta
// @Tags    apps
// @Produce json
// @Id      updateAppHealthcheckMeta
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Param   body body healthcheckdto.UpdateHealthcheckMetaReq true "request data"
// @Success 200 {object} healthcheckdto.UpdateHealthcheckMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/healthchecks/{itemID}/meta [put]
func (h *AppHandler) UpdateAppHealthcheckMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeHealthcheck, base.SettingScopeApp)
}

// DeleteAppHealthcheck Deletes healthcheck
// @Summary Deletes healthcheck
// @Description Deletes healthcheck
// @Tags    apps
// @Produce json
// @Id      deleteAppHealthcheck
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} healthcheckdto.DeleteHealthcheckResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/healthchecks/{itemID} [delete]
func (h *AppHandler) DeleteAppHealthcheck(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeHealthcheck, base.SettingScopeApp)
}

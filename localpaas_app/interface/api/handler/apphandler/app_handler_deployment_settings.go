package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// GetAppDeploymentSettings Gets app deployment settings
// @Summary Gets app deployment settings
// @Description Gets app deployment settings
// @Tags    apps
// @Produce json
// @Id      getAppDeploymentSettings
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.GetAppDeploymentSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/deployment-settings [get]
func (h *AppHandler) GetAppDeploymentSettings(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetAppDeploymentSettingsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetAppDeploymentSettings(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateAppDeploymentSettings Updates app deployment settings
// @Summary Updates app deployment settings
// @Description Updates app deployment settings
// @Tags    apps
// @Produce json
// @Id      updateAppDeploymentSettings
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.UpdateAppDeploymentSettingsReq true "request data"
// @Success 200 {object} appdto.UpdateAppDeploymentSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/deployment-settings [put]
func (h *AppHandler) UpdateAppDeploymentSettings(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewUpdateAppDeploymentSettingsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.UpdateAppDeploymentSettings(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

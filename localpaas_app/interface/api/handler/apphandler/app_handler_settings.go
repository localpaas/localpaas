package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// GetAppServiceSpec Gets app service spec
// @Summary Gets app service spec
// @Description Gets app service spec
// @Tags    apps
// @Produce json
// @Id      getAppServiceSpec
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.GetAppServiceSpecResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/service-spec [get]
func (h *AppHandler) GetAppServiceSpec(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetAppServiceSpecReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetAppServiceSpec(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

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
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeRead,
	})
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

// GetAppHttpSettings Gets app HTTP settings
// @Summary Gets app HTTP settings
// @Description Gets app HTTP settings
// @Tags    apps
// @Produce json
// @Id      getAppHttpSettings
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.GetAppHttpSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/http-settings [get]
func (h *AppHandler) GetAppHttpSettings(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetAppHttpSettingsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetAppHttpSettings(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetAppEnvVars Gets app Env vars
// @Summary Gets app Env vars
// @Description Gets app Env vars
// @Tags    apps
// @Produce json
// @Id      getAppEnvVars
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.GetAppEnvVarsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/env-vars [get]
func (h *AppHandler) GetAppEnvVars(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetAppEnvVarsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetAppEnvVars(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateAppServiceSpec Updates app service spec
// @Summary Updates app service spec
// @Description Updates app service spec
// @Tags    apps
// @Produce json
// @Id      updateAppServiceSpec
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.UpdateAppServiceSpecReq true "request data"
// @Success 200 {object} appdto.UpdateAppServiceSpecResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/service-spec [put]
func (h *AppHandler) UpdateAppServiceSpec(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewUpdateAppServiceSpecReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.UpdateAppServiceSpec(h.RequestCtx(ctx), auth, req)
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
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeWrite,
	})
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

// UpdateAppHttpSettings Updates app HTTP settings
// @Summary Updates app HTTP settings
// @Description Updates app HTTP settings
// @Tags    apps
// @Produce json
// @Id      updateAppHttpSettings
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.UpdateAppHttpSettingsReq true "request data"
// @Success 200 {object} appdto.UpdateAppHttpSettingsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/http-settings [put]
func (h *AppHandler) UpdateAppHttpSettings(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewUpdateAppHttpSettingsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.UpdateAppHttpSettings(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateAppEnvVars Updates app Env vars
// @Summary Updates app Env vars
// @Description Updates app Env vars
// @Tags    apps
// @Produce json
// @Id      updateAppEnvVars
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.UpdateAppEnvVarsReq true "request data"
// @Success 200 {object} appdto.UpdateAppEnvVarsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/env-vars [put]
func (h *AppHandler) UpdateAppEnvVars(ctx *gin.Context) {
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	appID, err := h.ParseStringParam(ctx, "appID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule:     base.ResourceModuleProject,
		ParentResourceType: base.ResourceTypeProject,
		ParentResourceID:   projectID,
		ResourceType:       base.ResourceTypeApp,
		ResourceID:         appID,
		Action:             base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewUpdateAppEnvVarsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.UpdateAppEnvVars(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// CreateApp Creates a new app
// @Summary Creates a new app
// @Description Creates a new app
// @Tags    apps
// @Produce json
// @Id      createApp
// @Param   projectID path string true "project ID"
// @Param   body body appdto.CreateAppReq true "request data"
// @Success 201 {object} appdto.CreateAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps [post]
func (h *AppHandler) CreateApp(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewCreateAppReq()
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.CreateApp(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateApp Updates an app
// @Summary Updates an app
// @Description Updates an app
// @Tags    apps
// @Produce json
// @Id      updateApp
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.UpdateAppReq true "request data"
// @Success 200 {object} appdto.UpdateAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID} [put]
func (h *AppHandler) UpdateApp(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewUpdateAppReq()
	req.ID = appID
	req.ProjectID = projectID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.UpdateApp(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteApp Deletes an app
// @Summary Deletes an app
// @Description Deletes an app
// @Tags    apps
// @Produce json
// @Id      deleteApp
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.DeleteAppReq true "request data"
// @Success 200 {object} appdto.DeleteAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID} [delete]
func (h *AppHandler) DeleteApp(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeDelete, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewDeleteAppReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.DeleteApp(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

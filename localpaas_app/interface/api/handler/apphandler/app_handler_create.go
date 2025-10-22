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
	projectID, err := h.ParseStringParam(ctx, "projectID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeProject,
		ResourceID:   projectID,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewCreateAppReq()
	req.ProjectID = projectID
	if err := h.ParseJSONBody(ctx, req); err != nil {
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
		ResourceType: base.ResourceTypeApp,
		ResourceID:   appID,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewDeleteAppReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseRequest(ctx, req, nil); err != nil {
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

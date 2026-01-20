package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// CreateAppTag Creates a new app tag
// @Summary Creates a new app tag
// @Description Creates a new app tag
// @Tags    apps
// @Produce json
// @Id      createAppTag
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.CreateAppTagReq true "request data"
// @Success 201 {object} appdto.CreateAppTagResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/tags [post]
func (h *AppHandler) CreateAppTag(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewCreateAppTagReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.CreateAppTag(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteAppTags Deletes app tags
// @Summary Deletes app tags
// @Description Deletes app tags
// @Tags    apps
// @Produce json
// @Id      deleteAppTag
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body appdto.DeleteAppTagsReq true "request data"
// @Success 200 {object} appdto.DeleteAppTagsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/tags/delete [post]
func (h *AppHandler) DeleteAppTags(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewDeleteAppTagsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.DeleteAppTags(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

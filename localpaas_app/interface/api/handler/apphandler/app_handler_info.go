package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// ListAppBase Lists apps
// @Summary Lists apps
// @Description Lists apps
// @Tags    apps
// @Produce json
// @Id      listAppBase
// @Param   projectID path string true "project ID"
// @Param   status query string false "`status=<target>`"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} appdto.ListAppBaseResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/base [get]
func (h *AppHandler) ListAppBase(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewListAppBaseReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.ListAppBase(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListApp Lists apps
// @Summary Lists apps
// @Description Lists apps
// @Tags    apps
// @Produce json
// @Id      listApp
// @Param   projectID path string true "project ID"
// @Param   status query string false "`status=<target>`"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} appdto.ListAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps [get]
func (h *AppHandler) ListApp(ctx *gin.Context) {
	auth, projectID, _, err := h.getAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewListAppReq()
	req.ProjectID = projectID
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.ListApp(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetApp Gets app details
// @Summary Gets app details
// @Description Gets app details
// @Tags    apps
// @Produce json
// @Id      getApp
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.GetAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID} [get]
func (h *AppHandler) GetApp(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetAppReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetApp(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

package apphandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/tiendc/gofn"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appdeploymentuc/appdeploymentdto"
)

// GetAppDeployment Gets app deployment
// @Summary Gets app deployment
// @Description Gets app deployment
// @Tags    apps_deployments
// @Produce json
// @Id      getAppDeployment
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "deployment ID"
// @Success 200 {object} appdeploymentdto.GetDeploymentResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/deployments/{id} [get]
func (h *AppHandler) GetAppDeployment(ctx *gin.Context) {
	auth, projectID, appID, itemID, err := h.getAuthForItem(ctx, base.ActionTypeRead)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdeploymentdto.NewGetDeploymentReq()
	req.ProjectID = projectID
	req.AppID = appID
	req.DeploymentID = itemID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appDeploymentUC.GetDeployment(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListAppDeployment Lists app deployments
// @Summary Lists app deployments
// @Description Lists app deployments
// @Tags    apps_deployments
// @Produce json
// @Id      listAppDeployment
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   status query string false "`status=<target>`"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} appdeploymentdto.ListDeploymentResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/deployments [get]
func (h *AppHandler) ListAppDeployment(ctx *gin.Context) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdeploymentdto.NewListDeploymentReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appDeploymentUC.ListDeployment(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetAppDeploymentLogs Stream app deployment logs via websocket
// @Summary Stream app deployment logs via websocket
// @Description Stream deployment app logs via websocket
// @Tags    apps_deployments
// @Produce json
// @Id      getAppDeploymentLogs
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "deployment ID"
// @Param   follow query string false "`follow=true/false`"
// @Param   since query string false "`since=YYYY-MM-DDTHH:mm:SSZ`"
// @Param   duration query int false "`duration=` logs within the period"
// @Param   tail query int false "`tail=1000` to get last 1000 lines of logs"
// @Success 200 {object} appdeploymentdto.GetDeploymentLogsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/deployments/{id}/logs [get]
func (h *AppHandler) GetAppDeploymentLogs(ctx *gin.Context, mel *melody.Melody) {
	auth, projectID, appID, itemID, err := h.getAuthForItem(ctx, base.ActionTypeRead)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdeploymentdto.NewGetDeploymentLogsReq()
	req.ProjectID = projectID
	req.AppID = appID
	req.DeploymentID = itemID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appDeploymentUC.GetDeploymentLogs(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Not a websocket request, return data via body
	if strings.ToLower(ctx.Request.Header.Get("Connection")) != "upgrade" {
		ctx.JSON(http.StatusOK, resp)
		return
	}

	go func() {
		for log := range resp.Data.LogChan {
			dataBytes := gofn.Must(json.Marshal(log))
			_ = mel.BroadcastBinaryFilter(dataBytes, func(session *melody.Session) bool {
				return session.Request == ctx.Request
			})
		}

		// Close the session
		for _, session := range gofn.Head(mel.Sessions()) {
			if session.Request == ctx.Request {
				_ = session.Close()
			}
		}
	}()

	_ = mel.HandleRequest(ctx.Writer, ctx.Request)
	_ = resp.Data.LogChanCloser()
}

// CancelAppDeployment Cancels app deployment
// @Summary Cancels app deployment
// @Description Cancels app deployment
// @Tags    apps_deployments
// @Produce json
// @Id      cancelAppDeployment
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "deployment ID"
// @Param   body body appdeploymentdto.CancelDeploymentReq true "request data"
// @Success 200 {object} appdeploymentdto.CancelDeploymentResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/deployments/{id}/cancel [post]
func (h *AppHandler) CancelAppDeployment(ctx *gin.Context) {
	auth, projectID, appID, itemID, err := h.getAuthForItem(ctx, base.ActionTypeWrite)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdeploymentdto.NewCancelDeploymentReq()
	req.ProjectID = projectID
	req.AppID = appID
	req.DeploymentID = itemID
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appDeploymentUC.CancelDeployment(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// GetAppRuntimeLogs Stream app logs via websocket
// @Summary Stream app logs via websocket
// @Description Stream app logs via websocket
// @Tags    apps
// @Produce json
// @Id      getAppRuntimeLogs
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   follow query string false "`follow=true/false`"
// @Param   since query string false "`since=YYYY-MM-DDTHH:mm:SSZ`"
// @Param   duration query int false "`duration=` logs within the period"
// @Param   tail query int false "`tail=1000` to get last 1000 lines of logs"
// @Success 200 {object} appdto.GetAppRuntimeLogsResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/runtime-logs [get]
func (h *AppHandler) GetAppRuntimeLogs(ctx *gin.Context, mel *melody.Melody) {
	auth, projectID, appID, err := h.getAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetAppRuntimeLogsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	isWebsocketReq := h.IsWebsocketRequest(ctx)
	if !isWebsocketReq {
		req.Follow = false // Not a websocket request, we don't support `follow` flag
	}

	resp, err := h.appUC.GetAppRuntimeLogs(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	if !isWebsocketReq {
		// Not a websocket request, return data via body
		ctx.JSON(http.StatusOK, resp)
	} else {
		h.StreamAppLogs(ctx, resp.Data.Logs, resp.Data.LogChan, resp.Data.LogChanCloser, mel)
	}
}

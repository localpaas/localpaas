package apphandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

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

	req := appdto.NewGetAppRuntimeLogsReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetAppRuntimeLogs(h.RequestCtx(ctx), auth, req)
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

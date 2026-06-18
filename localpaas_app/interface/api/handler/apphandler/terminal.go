package apphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
)

// GetAppTerminalInfo Gets terminal info
// @Summary Gets terminal info
// @Description Gets terminal info
// @Tags    apps
// @Produce json
// @Id      getAppTerminalInfo
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.GetTerminalInfoResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/terminal/info [get]
func (h *Handler) GetAppTerminalInfo(ctx *gin.Context) {
	auth, projectID, appID, err := h.GetAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewGetTerminalInfoReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.GetTerminalInfo(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// OpenAppTerminal Opens app terminal via websocket
// @Summary Opens app terminal via websocket
// @Description Opens app terminal via websocket
// @Tags    apps
// @Produce json
// @Id      getAppTerminal
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Success 200 {object} appdto.OpenTerminalResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/terminal [get]
func (h *Handler) OpenAppTerminal(ctx *gin.Context) {
	auth, projectID, appID, err := h.GetAuth(ctx, base.ActionTypeExecute, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := appdto.NewOpenTerminalReq()
	req.ProjectID = projectID
	req.AppID = appID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.appUC.OpenTerminal(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	if !h.IsWebsocketRequest(ctx) || resp.ExecAttachResult == nil {
		ctx.JSON(http.StatusOK, resp)
		return
	}
	defer resp.ExecAttachResult.Close()

	wsConn, err := h.UpgradeWebsocket(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	// Pipe container stdout/stderr (resp.ExecAttachResult.Reader) to websocket
	go func() {
		defer wsConn.Close()
		buf := make([]byte, 4096) //nolint:mnd
		for {
			n, err := resp.ExecAttachResult.Reader.Read(buf)
			if n > 0 {
				if err := wsConn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()

	// Pipe websocket messages to container stdin (resp.ExecAttachResult.Conn)
	for {
		mt, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		if mt == websocket.BinaryMessage || mt == websocket.TextMessage {
			if _, err := resp.ExecAttachResult.Conn.Write(message); err != nil {
				break
			}
		}
	}
}

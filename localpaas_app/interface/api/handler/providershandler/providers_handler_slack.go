package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
)

// ListSlack Lists Slack providers
// @Summary Lists Slack providers
// @Description Lists Slack providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderSlack
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} slackdto.ListSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack [get]
func (h *ProvidersHandler) ListSlack(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSlack, base.SettingScopeGlobal)
}

// GetSlack Gets Slack provider details
// @Summary Gets Slack provider details
// @Description Gets Slack provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderSlack
// @Param   id path string true "provider ID"
// @Success 200 {object} slackdto.GetSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{id} [get]
func (h *ProvidersHandler) GetSlack(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSlack, base.SettingScopeGlobal)
}

// CreateSlack Creates a new Slack provider
// @Summary Creates a new Slack provider
// @Description Creates a new Slack provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderSlack
// @Param   body body slackdto.CreateSlackReq true "request data"
// @Success 201 {object} slackdto.CreateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack [post]
func (h *ProvidersHandler) CreateSlack(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSlack, base.SettingScopeGlobal)
}

// UpdateSlack Updates Slack provider
// @Summary Updates Slack provider
// @Description Updates Slack provider
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSlack
// @Param   id path string true "provider ID"
// @Param   body body slackdto.UpdateSlackReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{id} [put]
func (h *ProvidersHandler) UpdateSlack(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSlack, base.SettingScopeGlobal)
}

// UpdateSlackMeta Updates Slack meta provider
// @Summary Updates Slack meta provider
// @Description Updates Slack meta provider
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSlackMeta
// @Param   id path string true "provider ID"
// @Param   body body slackdto.UpdateSlackMetaReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{id}/meta [put]
func (h *ProvidersHandler) UpdateSlackMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSlack, base.SettingScopeGlobal)
}

// DeleteSlack Deletes Slack provider
// @Summary Deletes Slack provider
// @Description Deletes Slack provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderSlack
// @Param   id path string true "provider ID"
// @Success 200 {object} slackdto.DeleteSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{id} [delete]
func (h *ProvidersHandler) DeleteSlack(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSlack, base.SettingScopeGlobal)
}

// TestSendSlackMsg Tests sending a msg
// @Summary Tests sending a msg
// @Description Tests sending a msg
// @Tags    global_providers
// @Produce json
// @Id      testSendSlackMsg
// @Param   body body slackdto.TestSendSlackMsgReq true "request data"
// @Success 200 {object} slackdto.TestSendSlackMsgResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/test-send-msg [post]
func (h *ProvidersHandler) TestSendSlackMsg(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewTestSendSlackMsgReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.SlackUC.TestSendSlackMsg(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

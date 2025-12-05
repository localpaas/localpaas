package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/slackuc/slackdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSlack Lists Slack providers
// @Summary Lists Slack providers
// @Description Lists Slack providers
// @Tags    providers_slack
// @Produce json
// @Id      listSlackProviders
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} slackdto.ListSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack [get]
func (h *ProvidersHandler) ListSlack(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSlack,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewListSlackReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.ListSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetSlack Gets Slack provider details
// @Summary Gets Slack provider details
// @Description Gets Slack provider details
// @Tags    providers_slack
// @Produce json
// @Id      getSlackProvider
// @Param   ID path string true "provider ID"
// @Success 200 {object} slackdto.GetSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{ID} [get]
func (h *ProvidersHandler) GetSlack(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSlack,
		ResourceID:     id,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewGetSlackReq()
	req.ID = id
	if err = h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.GetSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateSlack Creates a new Slack provider
// @Summary Creates a new Slack provider
// @Description Creates a new Slack provider
// @Tags    providers_slack
// @Produce json
// @Id      createSlackProvider
// @Param   body body slackdto.CreateSlackReq true "request data"
// @Success 201 {object} slackdto.CreateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack [post]
func (h *ProvidersHandler) CreateSlack(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewCreateSlackReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.CreateSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateSlack Updates Slack provider
// @Summary Updates Slack provider
// @Description Updates Slack provider
// @Tags    providers_slack
// @Produce json
// @Id      updateSlackProvider
// @Param   ID path string true "provider ID"
// @Param   body body slackdto.UpdateSlackReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{ID} [put]
func (h *ProvidersHandler) UpdateSlack(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSlack,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewUpdateSlackReq()
	req.ID = id
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.UpdateSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateSlackMeta Updates Slack meta provider
// @Summary Updates Slack meta provider
// @Description Updates Slack meta provider
// @Tags    providers_slack
// @Produce json
// @Id      updateSlackProviderMeta
// @Param   ID path string true "provider ID"
// @Param   body body slackdto.UpdateSlackMetaReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{ID}/meta [put]
func (h *ProvidersHandler) UpdateSlackMeta(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSlack,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewUpdateSlackMetaReq()
	req.ID = id
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.UpdateSlackMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteSlack Deletes Slack provider
// @Summary Deletes Slack provider
// @Description Deletes Slack provider
// @Tags    providers_slack
// @Produce json
// @Id      deleteSlackProvider
// @Param   ID path string true "provider ID"
// @Success 200 {object} slackdto.DeleteSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/slack/{ID} [delete]
func (h *ProvidersHandler) DeleteSlack(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleProvider,
		ResourceType:   base.ResourceTypeSlack,
		ResourceID:     id,
		Action:         base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := slackdto.NewDeleteSlackReq()
	req.ID = id
	if err := h.ParseRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.DeleteSlack(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// TestSendSlackMsg Tests sending a msg
// @Summary Tests sending a msg
// @Description Tests sending a msg
// @Tags    providers_slack
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
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.slackUC.TestSendSlackMsg(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

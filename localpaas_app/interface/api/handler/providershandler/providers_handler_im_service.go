package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/imserviceuc/imservicedto"
)

// ListIMService Lists IM services
// @Summary Lists IM services
// @Description Lists IM services
// @Tags    global_providers
// @Produce json
// @Id      listProviderIMService
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} imservicedto.ListIMServiceResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services [get]
func (h *ProvidersHandler) ListIMService(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeIMService, base.SettingScopeGlobal)
}

// GetIMService Gets IM service details
// @Summary Gets IM service details
// @Description Gets IM service details
// @Tags    global_providers
// @Produce json
// @Id      getProviderIMService
// @Param   itemID path string true "setting ID"
// @Success 200 {object} imservicedto.GetIMServiceResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services/{itemID} [get]
func (h *ProvidersHandler) GetIMService(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeIMService, base.SettingScopeGlobal)
}

// CreateIMService Creates a new IM service
// @Summary Creates a new IM service
// @Description Creates a new IM service
// @Tags    global_providers
// @Produce json
// @Id      createProviderIMService
// @Param   body body imservicedto.CreateIMServiceReq true "request data"
// @Success 201 {object} imservicedto.CreateIMServiceResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services [post]
func (h *ProvidersHandler) CreateIMService(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeIMService, base.SettingScopeGlobal)
}

// UpdateIMService Updates IM service
// @Summary Updates IM service
// @Description Updates IM service
// @Tags    global_providers
// @Produce json
// @Id      updateProviderIMService
// @Param   itemID path string true "setting ID"
// @Param   body body imservicedto.UpdateIMServiceReq true "request data"
// @Success 200 {object} imservicedto.UpdateIMServiceResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services/{itemID} [put]
func (h *ProvidersHandler) UpdateIMService(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeIMService, base.SettingScopeGlobal)
}

// UpdateIMServiceMeta Updates IMService meta provider
// @Summary Updates IMService meta provider
// @Description Updates IMService meta provider
// @Tags    global_providers
// @Produce json
// @Id      updateProviderIMServiceMeta
// @Param   itemID path string true "setting ID"
// @Param   body body imservicedto.UpdateIMServiceMetaReq true "request data"
// @Success 200 {object} imservicedto.UpdateIMServiceMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateIMServiceMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeIMService, base.SettingScopeGlobal)
}

// DeleteIMService Deletes IM service
// @Summary Deletes IM service
// @Description Deletes IM service
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderIMService
// @Param   itemID path string true "setting ID"
// @Success 200 {object} imservicedto.DeleteIMServiceResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services/{itemID} [delete]
func (h *ProvidersHandler) DeleteIMService(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeIMService, base.SettingScopeGlobal)
}

// TestSendInstantMsg Tests sending a msg
// @Summary Tests sending a msg
// @Description Tests sending a msg
// @Tags    global_providers
// @Produce json
// @Id      testSendIMServiceMsg
// @Param   body body imservicedto.TestSendInstantMsgReq true "request data"
// @Success 200 {object} imservicedto.TestSendInstantMsgResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/im-services/test-send-msg [post]
func (h *ProvidersHandler) TestSendInstantMsg(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := imservicedto.NewTestSendInstantMsgReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.IMServiceUC.TestSendInstantMsg(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

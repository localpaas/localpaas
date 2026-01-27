package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
)

// ListEmail Lists e-mail providers
// @Summary Lists e-mail providers
// @Description Lists e-mail providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderEmail
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} emaildto.ListEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails [get]
func (h *ProvidersHandler) ListEmail(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// GetEmail Gets e-mail provider details
// @Summary Gets e-mail provider details
// @Description Gets e-mail provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderEmail
// @Param   id path string true "provider ID"
// @Success 200 {object} emaildto.GetEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails/{id} [get]
func (h *ProvidersHandler) GetEmail(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// CreateEmail Creates a new e-mail provider
// @Summary Creates a new e-mail provider
// @Description Creates a new e-mail provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderEmail
// @Param   body body emaildto.CreateEmailReq true "request data"
// @Success 201 {object} emaildto.CreateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails [post]
func (h *ProvidersHandler) CreateEmail(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// UpdateEmail Updates e-mail provider
// @Summary Updates e-mail provider
// @Description Updates e-mail provider
// @Tags    global_providers
// @Produce json
// @Id      updateProviderEmail
// @Param   id path string true "provider ID"
// @Param   body body emaildto.UpdateEmailReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails/{id} [put]
func (h *ProvidersHandler) UpdateEmail(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// UpdateEmailMeta Updates email provider meta
// @Summary Updates email provider meta
// @Description Updates email provider meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderEmailMeta
// @Param   id path string true "provider ID"
// @Param   body body emaildto.UpdateEmailMetaReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails/{id}/meta [put]
func (h *ProvidersHandler) UpdateEmailMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// DeleteEmail Deletes e-mail provider
// @Summary Deletes e-mail provider
// @Description Deletes e-mail provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderEmail
// @Param   id path string true "provider ID"
// @Success 200 {object} emaildto.DeleteEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails/{id} [delete]
func (h *ProvidersHandler) DeleteEmail(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// TestSendEmail Tests sending an email
// @Summary Tests sending an email
// @Description Tests sending an email
// @Tags    global_providers
// @Produce json
// @Id      testSendEmail
// @Param   body body emaildto.TestSendEmailReq true "request data"
// @Success 200 {object} emaildto.TestSendEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/emails/test-send-mail [post]
func (h *ProvidersHandler) TestSendEmail(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := emaildto.NewTestSendEmailReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.EmailUC.TestSendEmail(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

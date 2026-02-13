package settinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
)

// ListEmail Lists e-mail settings
// @Summary Lists e-mail settings
// @Description Lists e-mail settings
// @Tags    settings
// @Produce json
// @Id      listSettingEmail
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} emaildto.ListEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails [get]
func (h *SettingHandler) ListEmail(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// GetEmail Gets e-mail setting details
// @Summary Gets e-mail setting details
// @Description Gets e-mail setting details
// @Tags    settings
// @Produce json
// @Id      getSettingEmail
// @Param   itemID path string true "setting ID"
// @Success 200 {object} emaildto.GetEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails/{itemID} [get]
func (h *SettingHandler) GetEmail(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// CreateEmail Creates a new e-mail setting
// @Summary Creates a new e-mail setting
// @Description Creates a new e-mail setting
// @Tags    settings
// @Produce json
// @Id      createSettingEmail
// @Param   body body emaildto.CreateEmailReq true "request data"
// @Success 201 {object} emaildto.CreateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails [post]
func (h *SettingHandler) CreateEmail(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// UpdateEmail Updates e-mail setting
// @Summary Updates e-mail setting
// @Description Updates e-mail setting
// @Tags    settings
// @Produce json
// @Id      updateSettingEmail
// @Param   itemID path string true "setting ID"
// @Param   body body emaildto.UpdateEmailReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails/{itemID} [put]
func (h *SettingHandler) UpdateEmail(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// UpdateEmailMeta Updates email setting meta
// @Summary Updates email setting meta
// @Description Updates email setting meta
// @Tags    settings
// @Produce json
// @Id      updateSettingEmailMeta
// @Param   itemID path string true "setting ID"
// @Param   body body emaildto.UpdateEmailMetaReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails/{itemID}/meta [put]
func (h *SettingHandler) UpdateEmailMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// DeleteEmail Deletes e-mail setting
// @Summary Deletes e-mail setting
// @Description Deletes e-mail setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingEmail
// @Param   itemID path string true "setting ID"
// @Success 200 {object} emaildto.DeleteEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails/{itemID} [delete]
func (h *SettingHandler) DeleteEmail(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeEmail, base.SettingScopeGlobal)
}

// TestSendMail Tests sending an email
// @Summary Tests sending an email
// @Description Tests sending an email
// @Tags    settings
// @Produce json
// @Id      testSendMail
// @Param   body body emaildto.TestSendMailReq true "request data"
// @Success 200 {object} emaildto.TestSendMailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/emails/test-send-mail [post]
func (h *SettingHandler) TestSendMail(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := emaildto.NewTestSendMailReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.EmailUC.TestSendMail(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

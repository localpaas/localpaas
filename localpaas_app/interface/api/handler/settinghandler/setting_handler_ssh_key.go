package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
)

// ListSSHKey Lists ssh-key settings
// @Summary Lists ssh-key settings
// @Description Lists ssh-key settings
// @Tags    settings
// @Produce json
// @Id      listSettingSSHKey
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} sshkeydto.ListSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssh-keys [get]
func (h *SettingHandler) ListSSHKey(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeGlobal)
}

// GetSSHKey Gets ssh-key setting details
// @Summary Gets ssh-key setting details
// @Description Gets ssh-key setting details
// @Tags    settings
// @Produce json
// @Id      getSettingSSHKey
// @Param   itemID path string true "setting ID"
// @Success 200 {object} sshkeydto.GetSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssh-keys/{itemID} [get]
func (h *SettingHandler) GetSSHKey(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeGlobal)
}

// CreateSSHKey Creates a new ssh-key setting
// @Summary Creates a new ssh-key setting
// @Description Creates a new ssh-key setting
// @Tags    settings
// @Produce json
// @Id      createSettingSSHKey
// @Param   body body sshkeydto.CreateSSHKeyReq true "request data"
// @Success 201 {object} sshkeydto.CreateSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssh-keys [post]
func (h *SettingHandler) CreateSSHKey(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeGlobal)
}

// UpdateSSHKey Updates ssh-key
// @Summary Updates ssh-key
// @Description Updates ssh-key
// @Tags    settings
// @Produce json
// @Id      updateSettingSSHKey
// @Param   itemID path string true "setting ID"
// @Param   body body sshkeydto.UpdateSSHKeyReq true "request data"
// @Success 200 {object} sshkeydto.UpdateSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssh-keys/{itemID} [put]
func (h *SettingHandler) UpdateSSHKey(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeGlobal)
}

// UpdateSSHKeyMeta Updates ssh-key meta
// @Summary Updates ssh-key meta
// @Description Updates ssh-key meta
// @Tags    settings
// @Produce json
// @Id      updateSettingSSHKeyMeta
// @Param   itemID path string true "setting ID"
// @Param   body body sshkeydto.UpdateSSHKeyMetaReq true "request data"
// @Success 200 {object} sshkeydto.UpdateSSHKeyMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssh-keys/{itemID}/meta [put]
func (h *SettingHandler) UpdateSSHKeyMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSHKey, base.SettingScopeGlobal)
}

// DeleteSSHKey Deletes sshkey setting
// @Summary Deletes sshkey setting
// @Description Deletes sshkey setting
// @Tags    settings
// @Produce json
// @Id      deleteSettingSSHKey
// @Param   itemID path string true "setting ID"
// @Success 200 {object} sshkeydto.DeleteSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/ssh-keys/{itemID} [delete]
func (h *SettingHandler) DeleteSSHKey(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeGlobal)
}

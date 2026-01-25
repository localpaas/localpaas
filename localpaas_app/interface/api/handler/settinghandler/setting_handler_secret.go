package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

// ListSecret Lists secrets
// @Summary Lists secrets
// @Description Lists secrets
// @Tags    settings
// @Produce json
// @Id      listSettingSecrets
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} secretdto.ListSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/secrets [get]
func (h *SettingHandler) ListSecret(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSecret, base.SettingScopeGlobal)
}

// CreateSecret Creates a new secret
// @Summary Creates a new secret
// @Description Creates a new secret
// @Tags    settings
// @Produce json
// @Id      createSettingSecret
// @Param   body body secretdto.CreateSecretReq true "request data"
// @Success 201 {object} secretdto.CreateSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/secrets [post]
func (h *SettingHandler) CreateSecret(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSecret, base.SettingScopeGlobal)
}

// UpdateSecret Updates secret
// @Summary Updates secret
// @Description Updates secret
// @Tags    settings
// @Produce json
// @Id      updateSettingSecret
// @Param   id path string true "setting ID"
// @Param   body body secretdto.UpdateSecretReq true "request data"
// @Success 201 {object} secretdto.UpdateSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/secrets/{id} [put]
func (h *SettingHandler) UpdateSecret(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSecret, base.SettingScopeGlobal)
}

// UpdateSecretMeta Updates secret meta
// @Summary Updates secret meta
// @Description Updates secret meta
// @Tags    settings
// @Produce json
// @Id      updateSettingSecretMeta
// @Param   id path string true "setting ID"
// @Param   body body secretdto.UpdateSecretMetaReq true "request data"
// @Success 201 {object} secretdto.UpdateSecretMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/secrets/{id}/meta [put]
func (h *SettingHandler) UpdateSecretMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSecret, base.SettingScopeGlobal)
}

// DeleteSecret Deletes a secret
// @Summary Deletes a secret
// @Description Deletes a secret
// @Tags    settings
// @Produce json
// @Id      deleteSettingSecret
// @Param   id path string true "setting ID"
// @Success 200 {object} secretdto.DeleteSecretResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/secrets/{id} [delete]
func (h *SettingHandler) DeleteSecret(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSecret, base.SettingScopeGlobal)
}

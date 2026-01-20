package settinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc/secretdto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

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
	auth, _, err := h.getAuth(ctx, base.ResourceTypeSecret, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewListSecretReq()
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.ListSecret(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, _, err := h.getAuth(ctx, base.ResourceTypeSecret, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewCreateSecretReq()
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.CreateSecret(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
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
	auth, id, err := h.getAuth(ctx, base.ResourceTypeSecret, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewUpdateSecretMetaReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.UpdateSecretMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, id, err := h.getAuth(ctx, base.ResourceTypeSecret, base.ActionTypeDelete, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := secretdto.NewDeleteSecretReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.secretUC.DeleteSecret(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

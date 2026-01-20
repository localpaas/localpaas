package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc/sshkeydto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSSHKey Lists ssh-key providers
// @Summary Lists ssh-key providers
// @Description Lists ssh-key providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderSSHKey
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} sshkeydto.ListSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssh-keys [get]
func (h *ProvidersHandler) ListSSHKey(ctx *gin.Context) {
	auth, _, err := h.getAuth(ctx, base.ResourceTypeSSHKey, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewListSSHKeyReq()
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.ListSSHKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetSSHKey Gets ssh-key provider details
// @Summary Gets ssh-key provider details
// @Description Gets ssh-key provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderSSHKey
// @Param   id path string true "provider ID"
// @Success 200 {object} sshkeydto.GetSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssh-keys/{id} [get]
func (h *ProvidersHandler) GetSSHKey(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeSSHKey, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewGetSSHKeyReq()
	req.ID = id
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.GetSSHKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateSSHKey Creates a new ssh-key provider
// @Summary Creates a new ssh-key provider
// @Description Creates a new ssh-key provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderSSHKey
// @Param   body body sshkeydto.CreateSSHKeyReq true "request data"
// @Success 201 {object} sshkeydto.CreateSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssh-keys [post]
func (h *ProvidersHandler) CreateSSHKey(ctx *gin.Context) {
	auth, _, err := h.getAuth(ctx, base.ResourceTypeSSHKey, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewCreateSSHKeyReq()
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.CreateSSHKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateSSHKey Updates ssh-key
// @Summary Updates ssh-key
// @Description Updates ssh-key
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSSHKey
// @Param   id path string true "provider ID"
// @Param   body body sshkeydto.UpdateSSHKeyReq true "request data"
// @Success 200 {object} sshkeydto.UpdateSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssh-keys/{id} [put]
func (h *ProvidersHandler) UpdateSSHKey(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeSSHKey, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewUpdateSSHKeyReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.UpdateSSHKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateSSHKeyMeta Updates ssh-key meta
// @Summary Updates ssh-key meta
// @Description Updates ssh-key meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderSSHKeyMeta
// @Param   id path string true "provider ID"
// @Param   body body sshkeydto.UpdateSSHKeyMetaReq true "request data"
// @Success 200 {object} sshkeydto.UpdateSSHKeyMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssh-keys/{id}/meta [put]
func (h *ProvidersHandler) UpdateSSHKeyMeta(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeSSHKey, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewUpdateSSHKeyMetaReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.UpdateSSHKeyMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteSSHKey Deletes sshkey provider
// @Summary Deletes sshkey provider
// @Description Deletes sshkey provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderSSHKey
// @Param   id path string true "provider ID"
// @Success 200 {object} sshkeydto.DeleteSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/ssh-keys/{id} [delete]
func (h *ProvidersHandler) DeleteSSHKey(ctx *gin.Context) {
	auth, id, err := h.getAuth(ctx, base.ResourceTypeSSHKey, base.ActionTypeDelete, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewDeleteSSHKeyReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.DeleteSSHKey(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

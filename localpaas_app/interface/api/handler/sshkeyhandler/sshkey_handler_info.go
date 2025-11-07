package sshkeyhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/sshkeyuc/sshkeydto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListSSHKeyBase Lists SSH key
// @Summary Lists SSH key
// @Description Lists SSH key
// @Tags    ssh_keys
// @Produce json
// @Id      listSSHKeyBase
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} sshkeydto.ListSSHKeyBaseResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /ssh-keys/base [get]
func (h *SSHKeyHandler) ListSSHKeyBase(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeSSHKey,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewListSSHKeyBaseReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.sshKeyUC.ListSSHKeyBase(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// ListSSHKey Lists SSH key
// @Summary Lists SSH key
// @Description Lists SSH key
// @Tags    ssh_keys
// @Produce json
// @Id      listSSHKey
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} sshkeydto.ListSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /ssh-keys [get]
func (h *SSHKeyHandler) ListSSHKey(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeSSHKey,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewListSSHKeyReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
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

// GetSSHKey Gets SSH key details
// @Summary Gets SSH key details
// @Description Gets SSH key details
// @Tags    ssh_keys
// @Produce json
// @Id      getSSHKey
// @Param   ID path string true "s3 storage ID"
// @Success 200 {object} sshkeydto.GetSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /ssh-keys/{ID} [get]
func (h *SSHKeyHandler) GetSSHKey(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeSSHKey,
		ResourceID:   id,
		Action:       base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := sshkeydto.NewGetSSHKeyReq()
	req.ID = id
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
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

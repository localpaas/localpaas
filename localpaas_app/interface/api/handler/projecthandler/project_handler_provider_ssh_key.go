package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
)

// ListSSHKey Lists ssh-key providers
// @Summary Lists ssh-key providers
// @Description Lists ssh-key providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectSSHKey
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} sshkeydto.ListSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssh-keys [get]
func (h *ProjectHandler) ListSSHKey(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeProject)
}

// GetSSHKey Gets ssh-key provider details
// @Summary Gets ssh-key provider details
// @Description Gets ssh-key provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectSSHKey
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} sshkeydto.GetSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssh-keys/{id} [get]
func (h *ProjectHandler) GetSSHKey(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeProject)
}

// CreateSSHKey Creates a new ssh-key provider
// @Summary Creates a new ssh-key provider
// @Description Creates a new ssh-key provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectSSHKey
// @Param   projectID path string true "project ID"
// @Param   body body sshkeydto.CreateSSHKeyReq true "request data"
// @Success 201 {object} sshkeydto.CreateSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssh-keys [post]
func (h *ProjectHandler) CreateSSHKey(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeProject)
}

// UpdateSSHKey Updates ssh-key
// @Summary Updates ssh-key
// @Description Updates ssh-key
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSSHKey
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body sshkeydto.UpdateSSHKeyReq true "request data"
// @Success 200 {object} sshkeydto.UpdateSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssh-keys/{id} [put]
func (h *ProjectHandler) UpdateSSHKey(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeProject)
}

// UpdateSSHKeyMeta Updates ssh-key meta
// @Summary Updates ssh-key meta
// @Description Updates ssh-key meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSSHKeyMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body sshkeydto.UpdateSSHKeyMetaReq true "request data"
// @Success 200 {object} sshkeydto.UpdateSSHKeyMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssh-keys/{id}/meta [put]
func (h *ProjectHandler) UpdateSSHKeyMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSHKey, base.SettingScopeProject)
}

// DeleteSSHKey Deletes ssh-key provider
// @Summary Deletes ssh-key provider
// @Description Deletes ssh-key provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectSSHKey
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} sshkeydto.DeleteSSHKeyResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/ssh-keys/{id} [delete]
func (h *ProjectHandler) DeleteSSHKey(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSHKey, base.SettingScopeProject)
}

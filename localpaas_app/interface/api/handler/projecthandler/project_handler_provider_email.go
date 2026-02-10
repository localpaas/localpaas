package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
)

// ListEmail Lists email providers
// @Summary Lists email providers
// @Description Lists email providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectEmail
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} emaildto.ListEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/emails [get]
func (h *ProjectHandler) ListEmail(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// GetEmail Gets email provider details
// @Summary Gets email provider details
// @Description Gets email provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectEmail
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} emaildto.GetEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/emails/{itemID} [get]
func (h *ProjectHandler) GetEmail(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// CreateEmail Creates a new email provider
// @Summary Creates a new email provider
// @Description Creates a new email provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectEmail
// @Param   projectID path string true "project ID"
// @Param   body body emaildto.CreateEmailReq true "request data"
// @Success 201 {object} emaildto.CreateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/emails [post]
func (h *ProjectHandler) CreateEmail(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// UpdateEmail Updates email provider
// @Summary Updates email provider
// @Description Updates email provider
// @Tags    project_providers
// @Produce json
// @Id      updateProjectEmail
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body emaildto.UpdateEmailReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/emails/{itemID} [put]
func (h *ProjectHandler) UpdateEmail(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// UpdateEmailMeta Updates Email meta provider
// @Summary Updates Email meta provider
// @Description Updates Email meta provider
// @Tags    project_providers
// @Produce json
// @Id      updateProjectEmailMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body emaildto.UpdateEmailMetaReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/emails/{itemID}/meta [put]
func (h *ProjectHandler) UpdateEmailMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// DeleteEmail Deletes email provider
// @Summary Deletes email provider
// @Description Deletes email provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectEmail
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} emaildto.DeleteEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/emails/{itemID} [delete]
func (h *ProjectHandler) DeleteEmail(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

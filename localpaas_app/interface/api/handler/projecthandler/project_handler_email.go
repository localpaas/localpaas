package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
)

// ListEmail Lists email settings
// @Summary Lists email settings
// @Description Lists email settings
// @Tags    project_settings
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
// @Router  /projects/{projectID}/emails [get]
func (h *ProjectHandler) ListEmail(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// GetEmail Gets email setting details
// @Summary Gets email setting details
// @Description Gets email setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectEmail
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} emaildto.GetEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/emails/{itemID} [get]
func (h *ProjectHandler) GetEmail(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// CreateEmail Creates a new email setting
// @Summary Creates a new email setting
// @Description Creates a new email setting
// @Tags    project_settings
// @Produce json
// @Id      createProjectEmail
// @Param   projectID path string true "project ID"
// @Param   body body emaildto.CreateEmailReq true "request data"
// @Success 201 {object} emaildto.CreateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/emails [post]
func (h *ProjectHandler) CreateEmail(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// UpdateEmail Updates email setting
// @Summary Updates email setting
// @Description Updates email setting
// @Tags    project_settings
// @Produce json
// @Id      updateProjectEmail
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body emaildto.UpdateEmailReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/emails/{itemID} [put]
func (h *ProjectHandler) UpdateEmail(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// UpdateEmailMeta Updates Email meta setting
// @Summary Updates Email meta setting
// @Description Updates Email meta setting
// @Tags    project_settings
// @Produce json
// @Id      updateProjectEmailMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body emaildto.UpdateEmailMetaReq true "request data"
// @Success 200 {object} emaildto.UpdateEmailMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/emails/{itemID}/meta [put]
func (h *ProjectHandler) UpdateEmailMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

// DeleteEmail Deletes email setting
// @Summary Deletes email setting
// @Description Deletes email setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectEmail
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} emaildto.DeleteEmailResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/emails/{itemID} [delete]
func (h *ProjectHandler) DeleteEmail(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeEmail, base.SettingScopeProject)
}

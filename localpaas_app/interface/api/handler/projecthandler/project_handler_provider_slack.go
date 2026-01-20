package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/slackuc/slackdto"
)

// ListSlack Lists Slack providers
// @Summary Lists Slack providers
// @Description Lists Slack providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectSlack
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} slackdto.ListSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack [get]
func (h *ProjectHandler) ListSlack(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSlack, base.SettingScopeProject)
}

// GetSlack Gets Slack provider details
// @Summary Gets Slack provider details
// @Description Gets Slack provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectSlack
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} slackdto.GetSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id} [get]
func (h *ProjectHandler) GetSlack(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSlack, base.SettingScopeProject)
}

// CreateSlack Creates a new Slack provider
// @Summary Creates a new Slack provider
// @Description Creates a new Slack provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectSlack
// @Param   projectID path string true "project ID"
// @Param   body body slackdto.CreateSlackReq true "request data"
// @Success 201 {object} slackdto.CreateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack [post]
func (h *ProjectHandler) CreateSlack(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSlack, base.SettingScopeProject)
}

// UpdateSlack Updates Slack provider
// @Summary Updates Slack provider
// @Description Updates Slack provider
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSlack
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body slackdto.UpdateSlackReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id} [put]
func (h *ProjectHandler) UpdateSlack(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSlack, base.SettingScopeProject)
}

// UpdateSlackMeta Updates Slack meta provider
// @Summary Updates Slack meta provider
// @Description Updates Slack meta provider
// @Tags    project_providers
// @Produce json
// @Id      updateProjectSlackMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body slackdto.UpdateSlackMetaReq true "request data"
// @Success 200 {object} slackdto.UpdateSlackMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id}/meta [put]
func (h *ProjectHandler) UpdateSlackMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSlack, base.SettingScopeProject)
}

// DeleteSlack Deletes Slack provider
// @Summary Deletes Slack provider
// @Description Deletes Slack provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectSlack
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} slackdto.DeleteSlackResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/slack/{id} [delete]
func (h *ProjectHandler) DeleteSlack(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSlack, base.SettingScopeProject)
}

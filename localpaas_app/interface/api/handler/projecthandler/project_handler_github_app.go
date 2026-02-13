package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
)

// ListGithubApp Lists github-app settings
// @Summary Lists github-app settings
// @Description Lists github-app settings
// @Tags    project_settings
// @Produce json
// @Id      listProjectGithubApp
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} githubappdto.ListGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps [get]
func (h *ProjectHandler) ListGithubApp(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeProject)
}

// GetGithubApp Gets github-app setting details
// @Summary Gets github-app setting details
// @Description Gets github-app setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectGithubApp
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} githubappdto.GetGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/{itemID} [get]
func (h *ProjectHandler) GetGithubApp(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeProject)
}

// CreateGithubApp Creates a new github-app setting
// @Summary Creates a new github-app setting
// @Description Creates a new github-app setting
// @Tags    project_settings
// @Produce json
// @Id      createProjectGithubApp
// @Param   projectID path string true "project ID"
// @Param   body body githubappdto.CreateGithubAppReq true "request data"
// @Success 201 {object} githubappdto.CreateGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps [post]
func (h *ProjectHandler) CreateGithubApp(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeProject)
}

// UpdateGithubApp Updates github-app
// @Summary Updates github-app
// @Description Updates github-app
// @Tags    project_settings
// @Produce json
// @Id      updateProjectGithubApp
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body githubappdto.UpdateGithubAppReq true "request data"
// @Success 200 {object} githubappdto.UpdateGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/{itemID} [put]
func (h *ProjectHandler) UpdateGithubApp(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeProject)
}

// UpdateGithubAppMeta Updates github-app meta
// @Summary Updates github-app meta
// @Description Updates github-app meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectGithubAppMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body githubappdto.UpdateGithubAppMetaReq true "request data"
// @Success 200 {object} githubappdto.UpdateGithubAppMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/{itemID}/meta [put]
func (h *ProjectHandler) UpdateGithubAppMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeGithubApp, base.SettingScopeProject)
}

// DeleteGithubApp Deletes github-app setting
// @Summary Deletes github-app setting
// @Description Deletes github-app setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectGithubApp
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} githubappdto.DeleteGithubAppResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/{itemID} [delete]
func (h *ProjectHandler) DeleteGithubApp(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeGithubApp, base.SettingScopeProject)
}

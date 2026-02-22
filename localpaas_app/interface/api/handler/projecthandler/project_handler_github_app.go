package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/githubappuc/githubappdto"
)

const (
	htmlSuccess = "<html><body><h3>Success</h3></body></html>"
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

// BeginProjectGithubAppManifestFlow Begins a github-app manifest flow
// @Summary Begins a github-app manifest flow
// @Description Begins a github-app manifest flow
// @Tags    project_settings
// @Produce json
// @Id      beginProjectGithubAppManifestFlow
// @Param   projectID path string true "project ID"
// @Param   body body githubappdto.BeginGithubAppManifestFlowReq true "request data"
// @Success 200 {object} githubappdto.BeginGithubAppManifestFlowResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/manifest-flow/begin [post]
func (h *ProjectHandler) BeginProjectGithubAppManifestFlow(ctx *gin.Context) {
	auth, projectID, _, err := h.GetAuthProjectSettings(ctx, base.ActionTypeWrite, "")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewBeginGithubAppManifestFlowReq()
	req.Scope = base.SettingScopeProject
	req.Type = base.SettingTypeGithubApp
	req.ObjectID = projectID
	if err = h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.BeginGithubAppManifestFlow(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// BeginProjectGithubAppManifestFlowCreation Begins a github-app manifest flow creation
// @Summary Begins a github-app manifest flow creation
// @Description Begins a github-app manifest flow creation
// @Tags    project_settings
// @Produce json
// @Id      beginProjectGithubAppManifestFlowCreation
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 "html page to redirect to github app creation page"
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/{itemID}/manifest-flow/begin [get]
func (h *ProjectHandler) BeginProjectGithubAppManifestFlowCreation(ctx *gin.Context) {
	itemID, err := h.ParseStringParam(ctx, "itemID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewBeginGithubAppManifestFlowCreationReq()
	req.SettingID = itemID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.BeginGithubAppManifestFlowCreation(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.Data(http.StatusOK, "text/html", reflectutil.UnsafeStrToBytes(resp.Data.PageContent))
}

// SetupProjectGithubAppManifestFlow Sets up a github-app manifest flow
// @Summary Sets up a github-app manifest flow
// @Description Sets up a github-app manifest flow
// @Tags    project_settings
// @Produce json
// @Id      setupProjectGithubAppManifestFlow
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 "html page to redirect to github app creation page"
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/github-apps/itemID}/manifest-flow/setup [get]
func (h *ProjectHandler) SetupProjectGithubAppManifestFlow(ctx *gin.Context) {
	itemID, err := h.ParseStringParam(ctx, "itemID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := githubappdto.NewSetupGithubAppManifestFlowReq()
	req.SettingID = itemID
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.GithubAppUC.SetupGithubAppManifestFlow(h.RequestCtx(ctx), req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	if resp.Data.RedirectURL != "" {
		ctx.Redirect(http.StatusFound, resp.Data.RedirectURL)
	}

	ctx.Data(http.StatusOK, "text/html", []byte(htmlSuccess))
}

package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/accesstokenuc/accesstokendto"
)

// ListAccessToken Lists access-token settings
// @Summary Lists access-token settings
// @Description Lists access-token settings
// @Tags    project_settings
// @Produce json
// @Id      listProjectAccessToken
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} accesstokendto.ListAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/access-tokens [get]
func (h *ProjectHandler) ListAccessToken(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeProject)
}

// GetAccessToken Gets access-token setting details
// @Summary Gets access-token setting details
// @Description Gets access-token setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectAccessToken
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} accesstokendto.GetAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/access-tokens/{itemID} [get]
func (h *ProjectHandler) GetAccessToken(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeProject)
}

// CreateAccessToken Creates a new access-token setting
// @Summary Creates a new access-token setting
// @Description Creates a new access-token setting
// @Tags    project_settings
// @Produce json
// @Id      createProjectAccessToken
// @Param   projectID path string true "project ID"
// @Param   body body accesstokendto.CreateAccessTokenReq true "request data"
// @Success 201 {object} accesstokendto.CreateAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/access-tokens [post]
func (h *ProjectHandler) CreateAccessToken(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeProject)
}

// UpdateAccessToken Updates access-token
// @Summary Updates access-token
// @Description Updates access-token
// @Tags    project_settings
// @Produce json
// @Id      updateProjectAccessToken
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body accesstokendto.UpdateAccessTokenReq true "request data"
// @Success 200 {object} accesstokendto.UpdateAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/access-tokens/{itemID} [put]
func (h *ProjectHandler) UpdateAccessToken(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeProject)
}

// UpdateAccessTokenMeta Updates access-token meta
// @Summary Updates access-token meta
// @Description Updates access-token meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectAccessTokenMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body accesstokendto.UpdateAccessTokenMetaReq true "request data"
// @Success 200 {object} accesstokendto.UpdateAccessTokenMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/access-tokens/{itemID}/meta [put]
func (h *ProjectHandler) UpdateAccessTokenMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeAccessToken, base.SettingScopeProject)
}

// DeleteAccessToken Deletes access-token setting
// @Summary Deletes access-token setting
// @Description Deletes access-token setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectAccessToken
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} accesstokendto.DeleteAccessTokenResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/access-tokens/{itemID} [delete]
func (h *ProjectHandler) DeleteAccessToken(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeAccessToken, base.SettingScopeProject)
}

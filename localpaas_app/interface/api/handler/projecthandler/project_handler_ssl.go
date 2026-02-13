package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/ssluc/ssldto"
)

// ListSSL Lists SSL settings
// @Summary Lists SSL settings
// @Description Lists SSL settings
// @Tags    project_settings
// @Produce json
// @Id      listProjectSSL
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} ssldto.ListSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssls [get]
func (h *ProjectHandler) ListSSL(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// GetSSL Gets SSL setting details
// @Summary Gets SSL setting details
// @Description Gets SSL setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectSSL
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} ssldto.GetSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssls/{itemID} [get]
func (h *ProjectHandler) GetSSL(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// CreateSSL Creates a new SSL setting
// @Summary Creates a new SSL setting
// @Description Creates a new SSL setting
// @Tags    project_settings
// @Produce json
// @Id      createProjectSSL
// @Param   projectID path string true "project ID"
// @Param   body body ssldto.CreateSSLReq true "request data"
// @Success 201 {object} ssldto.CreateSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssls [post]
func (h *ProjectHandler) CreateSSL(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// UpdateSSL Updates SSL
// @Summary Updates SSL
// @Description Updates SSL
// @Tags    project_settings
// @Produce json
// @Id      updateProjectSSL
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body ssldto.UpdateSSLReq true "request data"
// @Success 200 {object} ssldto.UpdateSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssls/{itemID} [put]
func (h *ProjectHandler) UpdateSSL(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// UpdateSSLMeta Updates SSL meta
// @Summary Updates SSL meta
// @Description Updates SSL meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectSSLMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body ssldto.UpdateSSLMetaReq true "request data"
// @Success 200 {object} ssldto.UpdateSSLMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssls/{itemID}/meta [put]
func (h *ProjectHandler) UpdateSSLMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

// DeleteSSL Deletes SSL setting
// @Summary Deletes SSL setting
// @Description Deletes SSL setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectSSL
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} ssldto.DeleteSSLResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssls/{itemID} [delete]
func (h *ProjectHandler) DeleteSSL(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSL, base.SettingScopeProject)
}

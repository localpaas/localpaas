package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
)

// ListAWS Lists AWS credentials
// @Summary Lists AWS credentials
// @Description Lists AWS credentials
// @Tags    project_settings
// @Produce json
// @Id      listProjectAWS
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} awsdto.ListAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws [get]
func (h *ProjectHandler) ListAWS(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeAWS, base.SettingScopeProject)
}

// GetAWS Gets AWS credential details
// @Summary Gets AWS credential details
// @Description Gets AWS credential details
// @Tags    project_settings
// @Produce json
// @Id      getProjectAWS
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awsdto.GetAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws/{itemID} [get]
func (h *ProjectHandler) GetAWS(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeAWS, base.SettingScopeProject)
}

// CreateAWS Creates a new AWS credential
// @Summary Creates a new AWS credential
// @Description Creates a new AWS credential
// @Tags    project_settings
// @Produce json
// @Id      createProjectAWS
// @Param   projectID path string true "project ID"
// @Param   body body awsdto.CreateAWSReq true "request data"
// @Success 201 {object} awsdto.CreateAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws [post]
func (h *ProjectHandler) CreateAWS(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeAWS, base.SettingScopeProject)
}

// UpdateAWS Updates AWS
// @Summary Updates AWS
// @Description Updates AWS
// @Tags    project_settings
// @Produce json
// @Id      updateProjectAWS
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body awsdto.UpdateAWSReq true "request data"
// @Success 200 {object} awsdto.UpdateAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws/{itemID} [put]
func (h *ProjectHandler) UpdateAWS(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeAWS, base.SettingScopeProject)
}

// UpdateAWSMeta Updates AWS credential meta
// @Summary Updates AWS credential meta
// @Description Updates AWS credential meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectAWSMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body awsdto.UpdateAWSMetaReq true "request data"
// @Success 200 {object} awsdto.UpdateAWSMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws/{itemID}/meta [put]
func (h *ProjectHandler) UpdateAWSMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeAWS, base.SettingScopeProject)
}

// DeleteAWS Deletes AWS credential
// @Summary Deletes AWS credential
// @Description Deletes AWS credential
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectAWS
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awsdto.DeleteAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws/{itemID} [delete]
func (h *ProjectHandler) DeleteAWS(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeAWS, base.SettingScopeProject)
}

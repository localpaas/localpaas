package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
)

// ListAWSS3 Lists S3 storage settings
// @Summary Lists S3 storage settings
// @Description Lists S3 storage settings
// @Tags    project_settings
// @Produce json
// @Id      listProjectAWSS3
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} awss3dto.ListAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws-s3 [get]
func (h *ProjectHandler) ListAWSS3(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeProject)
}

// GetAWSS3 Gets S3 storage setting details
// @Summary Gets S3 storage setting details
// @Description Gets S3 storage setting details
// @Tags    project_settings
// @Produce json
// @Id      getProjectAWSS3
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awss3dto.GetAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws-s3/{itemID} [get]
func (h *ProjectHandler) GetAWSS3(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeProject)
}

// CreateAWSS3 Creates a new S3 storage setting
// @Summary Creates a new S3 storage setting
// @Description Creates a new S3 storage setting
// @Tags    project_settings
// @Produce json
// @Id      createProjectAWSS3
// @Param   projectID path string true "project ID"
// @Param   body body awss3dto.CreateAWSS3Req true "request data"
// @Success 201 {object} awss3dto.CreateAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws-s3 [post]
func (h *ProjectHandler) CreateAWSS3(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeProject)
}

// UpdateAWSS3 Updates S3 storage
// @Summary Updates S3 storage
// @Description Updates S3 storage
// @Tags    project_settings
// @Produce json
// @Id      updateProjectAWSS3
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body awss3dto.UpdateAWSS3Req true "request data"
// @Success 200 {object} awss3dto.UpdateAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws-s3/{itemID} [put]
func (h *ProjectHandler) UpdateAWSS3(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeProject)
}

// UpdateAWSS3Meta Updates S3 storage meta
// @Summary Updates S3 storage meta
// @Description Updates S3 storage meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectAWSS3Meta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body awss3dto.UpdateAWSS3MetaReq true "request data"
// @Success 200 {object} awss3dto.UpdateAWSS3MetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws-s3/{itemID}/meta [put]
func (h *ProjectHandler) UpdateAWSS3Meta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeAWSS3, base.SettingScopeProject)
}

// DeleteAWSS3 Deletes S3 storage setting
// @Summary Deletes S3 storage setting
// @Description Deletes S3 storage setting
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectAWSS3
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awss3dto.DeleteAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/aws-s3/{itemID} [delete]
func (h *ProjectHandler) DeleteAWSS3(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeProject)
}

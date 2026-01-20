package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

// ListS3Storage Lists S3 storage providers
// @Summary Lists S3 storage providers
// @Description Lists S3 storage providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectS3Storage
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} s3storagedto.ListS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/s3-storages [get]
func (h *ProjectHandler) ListS3Storage(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeProject)
}

// GetS3Storage Gets S3 storage provider details
// @Summary Gets S3 storage provider details
// @Description Gets S3 storage provider details
// @Tags    project_providers
// @Produce json
// @Id      getProjectS3Storage
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} s3storagedto.GetS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/s3-storages/{id} [get]
func (h *ProjectHandler) GetS3Storage(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeProject)
}

// CreateS3Storage Creates a new S3 storage provider
// @Summary Creates a new S3 storage provider
// @Description Creates a new S3 storage provider
// @Tags    project_providers
// @Produce json
// @Id      createProjectS3Storage
// @Param   projectID path string true "project ID"
// @Param   body body s3storagedto.CreateS3StorageReq true "request data"
// @Success 201 {object} s3storagedto.CreateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/s3-storages [post]
func (h *ProjectHandler) CreateS3Storage(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeProject)
}

// UpdateS3Storage Updates S3 storage
// @Summary Updates S3 storage
// @Description Updates S3 storage
// @Tags    project_providers
// @Produce json
// @Id      updateProjectS3Storage
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body s3storagedto.UpdateS3StorageReq true "request data"
// @Success 200 {object} s3storagedto.UpdateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/s3-storages/{id} [put]
func (h *ProjectHandler) UpdateS3Storage(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeProject)
}

// UpdateS3StorageMeta Updates S3 storage meta
// @Summary Updates S3 storage meta
// @Description Updates S3 storage meta
// @Tags    project_providers
// @Produce json
// @Id      updateProjectS3StorageMeta
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Param   body body s3storagedto.UpdateS3StorageMetaReq true "request data"
// @Success 200 {object} s3storagedto.UpdateS3StorageMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/s3-storages/{id}/meta [put]
func (h *ProjectHandler) UpdateS3StorageMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeS3Storage, base.SettingScopeProject)
}

// DeleteS3Storage Deletes S3 storage provider
// @Summary Deletes S3 storage provider
// @Description Deletes S3 storage provider
// @Tags    project_providers
// @Produce json
// @Id      deleteProjectS3Storage
// @Param   projectID path string true "project ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} s3storagedto.DeleteS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/providers/s3-storages/{id} [delete]
func (h *ProjectHandler) DeleteS3Storage(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeProject)
}

package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

// ListS3Storage Lists S3 storage providers
// @Summary Lists S3 storage providers
// @Description Lists S3 storage providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderS3Storage
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} s3storagedto.ListS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages [get]
func (h *ProvidersHandler) ListS3Storage(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeGlobal)
}

// GetS3Storage Gets S3 storage provider details
// @Summary Gets S3 storage provider details
// @Description Gets S3 storage provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderS3Storage
// @Param   id path string true "provider ID"
// @Success 200 {object} s3storagedto.GetS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages/{id} [get]
func (h *ProvidersHandler) GetS3Storage(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeGlobal)
}

// CreateS3Storage Creates a new S3 storage provider
// @Summary Creates a new S3 storage provider
// @Description Creates a new S3 storage provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderS3Storage
// @Param   body body s3storagedto.CreateS3StorageReq true "request data"
// @Success 201 {object} s3storagedto.CreateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages [post]
func (h *ProvidersHandler) CreateS3Storage(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeGlobal)
}

// UpdateS3Storage Updates S3 storage
// @Summary Updates S3 storage
// @Description Updates S3 storage
// @Tags    global_providers
// @Produce json
// @Id      updateProviderS3Storage
// @Param   id path string true "provider ID"
// @Param   body body s3storagedto.UpdateS3StorageReq true "request data"
// @Success 200 {object} s3storagedto.UpdateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages/{id} [put]
func (h *ProvidersHandler) UpdateS3Storage(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeGlobal)
}

// UpdateS3StorageMeta Updates S3 storage meta
// @Summary Updates S3 storage meta
// @Description Updates S3 storage meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderS3StorageMeta
// @Param   id path string true "provider ID"
// @Param   body body s3storagedto.UpdateS3StorageMetaReq true "request data"
// @Success 200 {object} s3storagedto.UpdateS3StorageMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages/{id}/meta [put]
func (h *ProvidersHandler) UpdateS3StorageMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeS3Storage, base.SettingScopeGlobal)
}

// DeleteS3Storage Deletes S3 storage provider
// @Summary Deletes S3 storage provider
// @Description Deletes S3 storage provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderS3Storage
// @Param   id path string true "provider ID"
// @Success 200 {object} s3storagedto.DeleteS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages/{id} [delete]
func (h *ProvidersHandler) DeleteS3Storage(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeS3Storage, base.SettingScopeGlobal)
}

// TestS3StorageConn Test S3 storage connection
// @Summary Test S3 storage connection
// @Description Test S3 storage connection
// @Tags    global_providers
// @Produce json
// @Id      testS3StorageConn
// @Param   body body s3storagedto.TestS3StorageConnReq true "request data"
// @Success 200 {object} s3storagedto.TestS3StorageConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/s3-storages/test-conn [post]
func (h *ProvidersHandler) TestS3StorageConn(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewTestS3StorageConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.S3StorageUC.TestS3StorageConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc/s3storagedto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

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
	auth, _, err := h.getAuth(ctx, base.ResourceTypeS3Storage, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewListS3StorageReq()
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.ListS3Storage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, id, err := h.getAuth(ctx, base.ResourceTypeS3Storage, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewGetS3StorageReq()
	req.ID = id
	req.GlobalOnly = true
	if err = h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.GetS3Storage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, _, err := h.getAuth(ctx, base.ResourceTypeS3Storage, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewCreateS3StorageReq()
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.CreateS3Storage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
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
	auth, id, err := h.getAuth(ctx, base.ResourceTypeS3Storage, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewUpdateS3StorageReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.UpdateS3Storage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, id, err := h.getAuth(ctx, base.ResourceTypeS3Storage, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewUpdateS3StorageMetaReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.UpdateS3StorageMeta(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, id, err := h.getAuth(ctx, base.ResourceTypeS3Storage, base.ActionTypeDelete, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewDeleteS3StorageReq()
	req.ID = id
	req.GlobalOnly = true
	if err := h.ParseAndValidateRequest(ctx, req, nil); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.DeleteS3Storage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
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
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewTestS3StorageConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.s3StorageUC.TestS3StorageConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

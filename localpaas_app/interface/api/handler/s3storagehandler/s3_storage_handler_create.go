package s3storagehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/s3storageuc/s3storagedto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// CreateS3Storage Creates a new s3 storage
// @Summary Creates a new s3 storage
// @Description Creates a new s3 storage
// @Tags    s3_storages
// @Produce json
// @Id      createS3Storage
// @Param   body body s3storagedto.CreateS3StorageReq true "request data"
// @Success 201 {object} s3storagedto.CreateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /s3-storages [post]
func (h *S3StorageHandler) CreateS3Storage(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeS3Storage,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewCreateS3StorageReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
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

// UpdateS3Storage Updates a s3 storage
// @Summary Updates a s3 storage
// @Description Updates a s3 storage
// @Tags    s3_storages
// @Produce json
// @Id      updateS3Storage
// @Param   ID path string true "s3 storage ID"
// @Success 200 {object} s3storagedto.UpdateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /s3-storages/{ID} [put]
func (h *S3StorageHandler) UpdateS3Storage(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeS3Storage,
		ResourceID:   id,
		Action:       base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewUpdateS3StorageReq()
	req.ID = id
	if err := h.ParseJSONBody(ctx, req); err != nil {
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

// DeleteS3Storage Deletes a s3 storage
// @Summary Deletes a s3 storage
// @Description Deletes a s3 storage
// @Tags    s3_storages
// @Produce json
// @Id      deleteS3Storage
// @Param   ID path string true "s3 storage ID"
// @Success 200 {object} s3storagedto.DeleteS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /s3-storages/{ID} [delete]
func (h *S3StorageHandler) DeleteS3Storage(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceType: base.ResourceTypeS3Storage,
		ResourceID:   id,
		Action:       base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewDeleteS3StorageReq()
	req.ID = id
	if err := h.ParseRequest(ctx, req, nil); err != nil {
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

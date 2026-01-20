package projecthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/s3storageuc/s3storagedto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListS3Storage Lists S3 storage providers
// @Summary Lists S3 storage providers
// @Description Lists S3 storage providers
// @Tags    project_providers
// @Produce json
// @Id      listProjectS3Storages
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
	auth, projectID, _, err := h.getProjectProviderAuth(ctx, base.ActionTypeRead, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewListS3StorageReq()
	req.ProjectID = projectID
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
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeRead, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewGetS3StorageReq()
	req.ID = id
	req.ProjectID = projectID
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
	auth, projectID, _, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, false)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewCreateS3StorageReq()
	req.ProjectID = projectID
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
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewUpdateS3StorageReq()
	req.ID = id
	req.ProjectID = projectID
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
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewUpdateS3StorageMetaReq()
	req.ID = id
	req.ProjectID = projectID
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
	auth, projectID, id, err := h.getProjectProviderAuth(ctx, base.ActionTypeWrite, true)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewDeleteS3StorageReq()
	req.ID = id
	req.ProjectID = projectID
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

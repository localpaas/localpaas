package settingshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/s3storageuc/s3storagedto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListS3Storage Lists S3 storage settings
// @Summary Lists S3 storage settings
// @Description Lists S3 storage settings
// @Tags    settings_s3_storage
// @Produce json
// @Id      listS3StorageSettings
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} s3storagedto.ListS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages [get]
func (h *SettingsHandler) ListS3Storage(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeS3Storage,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewListS3StorageReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
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

// GetS3Storage Gets S3 storage setting details
// @Summary Gets S3 storage setting details
// @Description Gets S3 storage setting details
// @Tags    settings_s3_storage
// @Produce json
// @Id      getS3StorageSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} s3storagedto.GetS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages/{ID} [get]
func (h *SettingsHandler) GetS3Storage(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeS3Storage,
		ResourceID:     id,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewGetS3StorageReq()
	req.ID = id
	if err = h.ParseRequest(ctx, req, nil); err != nil {
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

// CreateS3Storage Creates a new S3 storage setting
// @Summary Creates a new S3 storage setting
// @Description Creates a new S3 storage setting
// @Tags    settings_s3_storage
// @Produce json
// @Id      createS3StorageSetting
// @Param   body body s3storagedto.CreateS3StorageReq true "request data"
// @Success 201 {object} s3storagedto.CreateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages [post]
func (h *SettingsHandler) CreateS3Storage(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		Action:         base.ActionTypeWrite,
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

// UpdateS3Storage Updates S3 storage
// @Summary Updates S3 storage
// @Description Updates S3 storage
// @Tags    settings_s3_storage
// @Produce json
// @Id      updateS3StorageSetting
// @Param   ID path string true "setting ID"
// @Param   body body s3storagedto.UpdateS3StorageReq true "request data"
// @Success 200 {object} s3storagedto.UpdateS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages/{ID} [put]
func (h *SettingsHandler) UpdateS3Storage(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeS3Storage,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
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

// UpdateS3StorageMeta Updates S3 storage meta
// @Summary Updates S3 storage meta
// @Description Updates S3 storage meta
// @Tags    settings_s3_storage
// @Produce json
// @Id      updateS3StorageMetaSetting
// @Param   ID path string true "setting ID"
// @Param   body body s3storagedto.UpdateS3StorageMetaReq true "request data"
// @Success 200 {object} s3storagedto.UpdateS3StorageMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages/{ID}/meta [put]
func (h *SettingsHandler) UpdateS3StorageMeta(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeS3Storage,
		ResourceID:     id,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewUpdateS3StorageMetaReq()
	req.ID = id
	if err := h.ParseJSONBody(ctx, req); err != nil {
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

// DeleteS3Storage Deletes S3 storage setting
// @Summary Deletes S3 storage setting
// @Description Deletes S3 storage setting
// @Tags    settings_s3_storage
// @Produce json
// @Id      deleteS3StorageSetting
// @Param   ID path string true "setting ID"
// @Success 200 {object} s3storagedto.DeleteS3StorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages/{ID} [delete]
func (h *SettingsHandler) DeleteS3Storage(ctx *gin.Context) {
	id, err := h.ParseStringParam(ctx, "ID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleSettings,
		ResourceType:   base.ResourceTypeS3Storage,
		ResourceID:     id,
		Action:         base.ActionTypeDelete,
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

// TestS3StorageConn Test S3 storage connection
// @Summary Test S3 storage connection
// @Description Test S3 storage connection
// @Tags    settings_s3_storage
// @Produce json
// @Id      testS3StorageConn
// @Param   body body s3storagedto.TestS3StorageConnReq true "request data"
// @Success 200 {object} s3storagedto.TestS3StorageConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/s3-storages/test-conn [post]
func (h *SettingsHandler) TestS3StorageConn(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := s3storagedto.NewTestS3StorageConnReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
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

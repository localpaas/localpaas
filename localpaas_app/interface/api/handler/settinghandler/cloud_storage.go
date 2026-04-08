package settinghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cloudstorageuc/cloudstoragedto"
)

// ListCloudStorage Lists cloud storages
// @Summary Lists cloud storages
// @Description Lists cloud storages
// @Tags    settings
// @Produce json
// @Id      listSettingCloudStorage
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} cloudstoragedto.ListCloudStorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages [get]
func (h *Handler) ListCloudStorage(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeCloudStorage, base.SettingScopeGlobal)
}

// GetCloudStorage Gets cloud storage details
// @Summary Gets cloud storage details
// @Description Gets cloud storage details
// @Tags    settings
// @Produce json
// @Id      getSettingCloudStorage
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cloudstoragedto.GetCloudStorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages/{itemID} [get]
func (h *Handler) GetCloudStorage(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeCloudStorage, base.SettingScopeGlobal)
}

// CreateCloudStorage Creates a new cloud storage
// @Summary Creates a new cloud storage
// @Description Creates a new cloud storage
// @Tags    settings
// @Produce json
// @Id      createSettingCloudStorage
// @Param   body body cloudstoragedto.CreateCloudStorageReq true "request data"
// @Success 201 {object} cloudstoragedto.CreateCloudStorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages [post]
func (h *Handler) CreateCloudStorage(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeCloudStorage, base.SettingScopeGlobal)
}

// UpdateCloudStorage Updates a cloud storage
// @Summary Updates a cloud storage
// @Description Updates a cloud storage
// @Tags    settings
// @Produce json
// @Id      updateSettingCloudStorage
// @Param   itemID path string true "setting ID"
// @Param   body body cloudstoragedto.UpdateCloudStorageReq true "request data"
// @Success 200 {object} cloudstoragedto.UpdateCloudStorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages/{itemID} [put]
func (h *Handler) UpdateCloudStorage(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeCloudStorage, base.SettingScopeGlobal)
}

// UpdateCloudStorageMeta Updates a cloud storage's meta
// @Summary Updates a cloud storage's meta
// @Description Updates a cloud storage's meta
// @Tags    settings
// @Produce json
// @Id      updateSettingCloudStorageMeta
// @Param   itemID path string true "setting ID"
// @Param   body body cloudstoragedto.UpdateCloudStorageMetaReq true "request data"
// @Success 200 {object} cloudstoragedto.UpdateCloudStorageMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages/{itemID}/meta [put]
func (h *Handler) UpdateCloudStorageMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeCloudStorage, base.SettingScopeGlobal)
}

// DeleteCloudStorage Deletes a cloud storage
// @Summary Deletes a cloud storage
// @Description Deletes a cloud storage
// @Tags    settings
// @Produce json
// @Id      deleteSettingCloudStorage
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cloudstoragedto.DeleteCloudStorageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages/{itemID} [delete]
func (h *Handler) DeleteCloudStorage(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeCloudStorage, base.SettingScopeGlobal)
}

// TestCloudStorageConn Test cloud storage connection
// @Summary Test cloud storage connection
// @Description Test cloud storage connection
// @Tags    settings
// @Produce json
// @Id      testCloudStorageConn
// @Param   body body cloudstoragedto.TestCloudStorageConnReq true "request data"
// @Success 200 {object} cloudstoragedto.TestCloudStorageConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cloud-storages/test-conn [post]
func (h *Handler) TestCloudStorageConn(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := cloudstoragedto.NewTestCloudStorageConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.CloudStorageUC.TestCloudStorageConn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

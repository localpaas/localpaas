package providershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/awss3uc/awss3dto"
)

// ListAWSS3 Lists S3 storage providers
// @Summary Lists S3 storage providers
// @Description Lists S3 storage providers
// @Tags    global_providers
// @Produce json
// @Id      listProviderAWSS3
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} awss3dto.ListAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3 [get]
func (h *ProvidersHandler) ListAWSS3(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeGlobal)
}

// GetAWSS3 Gets S3 storage provider details
// @Summary Gets S3 storage provider details
// @Description Gets S3 storage provider details
// @Tags    global_providers
// @Produce json
// @Id      getProviderAWSS3
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awss3dto.GetAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3/{itemID} [get]
func (h *ProvidersHandler) GetAWSS3(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeGlobal)
}

// CreateAWSS3 Creates a new S3 storage provider
// @Summary Creates a new S3 storage provider
// @Description Creates a new S3 storage provider
// @Tags    global_providers
// @Produce json
// @Id      createProviderAWSS3
// @Param   body body awss3dto.CreateAWSS3Req true "request data"
// @Success 201 {object} awss3dto.CreateAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3 [post]
func (h *ProvidersHandler) CreateAWSS3(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeGlobal)
}

// UpdateAWSS3 Updates S3 storage
// @Summary Updates S3 storage
// @Description Updates S3 storage
// @Tags    global_providers
// @Produce json
// @Id      updateProviderAWSS3
// @Param   itemID path string true "setting ID"
// @Param   body body awss3dto.UpdateAWSS3Req true "request data"
// @Success 200 {object} awss3dto.UpdateAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3/{itemID} [put]
func (h *ProvidersHandler) UpdateAWSS3(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeGlobal)
}

// UpdateAWSS3Meta Updates S3 storage meta
// @Summary Updates S3 storage meta
// @Description Updates S3 storage meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderAWSS3Meta
// @Param   itemID path string true "setting ID"
// @Param   body body awss3dto.UpdateAWSS3MetaReq true "request data"
// @Success 200 {object} awss3dto.UpdateAWSS3MetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateAWSS3Meta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeAWSS3, base.SettingScopeGlobal)
}

// DeleteAWSS3 Deletes S3 storage provider
// @Summary Deletes S3 storage provider
// @Description Deletes S3 storage provider
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderAWSS3
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awss3dto.DeleteAWSS3Resp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3/{itemID} [delete]
func (h *ProvidersHandler) DeleteAWSS3(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeAWSS3, base.SettingScopeGlobal)
}

// TestAWSS3Conn Test S3 storage connection
// @Summary Test S3 storage connection
// @Description Test S3 storage connection
// @Tags    global_providers
// @Produce json
// @Id      testAWSS3Conn
// @Param   body body awss3dto.TestAWSS3ConnReq true "request data"
// @Success 200 {object} awss3dto.TestAWSS3ConnResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws-s3/test-conn [post]
func (h *ProvidersHandler) TestAWSS3Conn(ctx *gin.Context) {
	auth, err := h.AuthHandler.GetCurrentAuth(ctx, authhandler.NoAccessCheck)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := awss3dto.NewTestAWSS3ConnReq()
	if err := h.ParseAndValidateJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.AWSS3UC.TestAWSS3Conn(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

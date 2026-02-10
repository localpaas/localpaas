package providershandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/awsuc/awsdto"
)

// ListAWS Lists AWS credentials
// @Summary Lists AWS credentials
// @Description Lists AWS credentials
// @Tags    global_providers
// @Produce json
// @Id      listProviderAWS
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} awsdto.ListAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws [get]
func (h *ProvidersHandler) ListAWS(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeAWS, base.SettingScopeGlobal)
}

// GetAWS Gets AWS credential details
// @Summary Gets AWS credential details
// @Description Gets AWS credential details
// @Tags    global_providers
// @Produce json
// @Id      getProviderAWS
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awsdto.GetAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws/{itemID} [get]
func (h *ProvidersHandler) GetAWS(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeAWS, base.SettingScopeGlobal)
}

// CreateAWS Creates a new AWS credential
// @Summary Creates a new AWS credential
// @Description Creates a new AWS credential
// @Tags    global_providers
// @Produce json
// @Id      createProviderAWS
// @Param   body body awsdto.CreateAWSReq true "request data"
// @Success 201 {object} awsdto.CreateAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws [post]
func (h *ProvidersHandler) CreateAWS(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeAWS, base.SettingScopeGlobal)
}

// UpdateAWS Updates S3 storage
// @Summary Updates S3 storage
// @Description Updates S3 storage
// @Tags    global_providers
// @Produce json
// @Id      updateProviderAWS
// @Param   itemID path string true "setting ID"
// @Param   body body awsdto.UpdateAWSReq true "request data"
// @Success 200 {object} awsdto.UpdateAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws/{itemID} [put]
func (h *ProvidersHandler) UpdateAWS(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeAWS, base.SettingScopeGlobal)
}

// UpdateAWSMeta Updates S3 storage meta
// @Summary Updates S3 storage meta
// @Description Updates S3 storage meta
// @Tags    global_providers
// @Produce json
// @Id      updateProviderAWSMeta
// @Param   itemID path string true "setting ID"
// @Param   body body awsdto.UpdateAWSMetaReq true "request data"
// @Success 200 {object} awsdto.UpdateAWSMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws/{itemID}/meta [put]
func (h *ProvidersHandler) UpdateAWSMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeAWS, base.SettingScopeGlobal)
}

// DeleteAWS Deletes AWS credential
// @Summary Deletes AWS credential
// @Description Deletes AWS credential
// @Tags    global_providers
// @Produce json
// @Id      deleteProviderAWS
// @Param   itemID path string true "setting ID"
// @Success 200 {object} awsdto.DeleteAWSResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /providers/aws/{itemID} [delete]
func (h *ProvidersHandler) DeleteAWS(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeAWS, base.SettingScopeGlobal)
}

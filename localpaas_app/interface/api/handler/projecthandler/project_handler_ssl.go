package projecthandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/sslcertuc/sslcertdto"
)

// ListSSLCert Lists SSL certs
// @Summary Lists SSL certs
// @Description Lists SSL certs
// @Tags    project_settings
// @Produce json
// @Id      listProjectSSLCert
// @Param   projectID path string true "project ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} sslcertdto.ListSSLCertResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-certs [get]
func (h *ProjectHandler) ListSSLCert(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeSSLCert, base.SettingScopeProject)
}

// GetSSLCert Gets SSL cert details
// @Summary Gets SSL cert details
// @Description Gets SSL cert details
// @Tags    project_settings
// @Produce json
// @Id      getProjectSSLCert
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} sslcertdto.GetSSLCertResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-certs/{itemID} [get]
func (h *ProjectHandler) GetSSLCert(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeSSLCert, base.SettingScopeProject)
}

// CreateSSLCert Creates a new SSL cert
// @Summary Creates a new SSL cert
// @Description Creates a new SSL cert
// @Tags    project_settings
// @Produce json
// @Id      createProjectSSLCert
// @Param   projectID path string true "project ID"
// @Param   body body sslcertdto.CreateSSLCertReq true "request data"
// @Success 201 {object} sslcertdto.CreateSSLCertResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-certs [post]
func (h *ProjectHandler) CreateSSLCert(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeSSLCert, base.SettingScopeProject)
}

// UpdateSSLCert Updates an SSL cert
// @Summary Updates an SSL cert
// @Description Updates an SSL cert
// @Tags    project_settings
// @Produce json
// @Id      updateProjectSSLCert
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body sslcertdto.UpdateSSLCertReq true "request data"
// @Success 200 {object} sslcertdto.UpdateSSLCertResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-certs/{itemID} [put]
func (h *ProjectHandler) UpdateSSLCert(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeSSLCert, base.SettingScopeProject)
}

// UpdateSSLCertMeta Updates SSL cert meta
// @Summary Updates SSL cert meta
// @Description Updates SSL cert meta
// @Tags    project_settings
// @Produce json
// @Id      updateProjectSSLCertMeta
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Param   body body sslcertdto.UpdateSSLCertMetaReq true "request data"
// @Success 200 {object} sslcertdto.UpdateSSLCertMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-certs/{itemID}/meta [put]
func (h *ProjectHandler) UpdateSSLCertMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeSSLCert, base.SettingScopeProject)
}

// DeleteSSLCert Deletes an SSL cert
// @Summary Deletes an SSL cert
// @Description Deletes an SSL cert
// @Tags    project_settings
// @Produce json
// @Id      deleteProjectSSLCert
// @Param   projectID path string true "project ID"
// @Param   itemID path string true "setting ID"
// @Success 200 {object} sslcertdto.DeleteSSLCertResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/ssl-certs/{itemID} [delete]
func (h *ProjectHandler) DeleteSSLCert(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeSSLCert, base.SettingScopeProject)
}

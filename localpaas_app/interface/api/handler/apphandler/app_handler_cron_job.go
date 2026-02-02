package apphandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

// ListAppCronJob Lists cron-jobs
// @Summary Lists cron-jobs
// @Description Lists cron-jobs
// @Tags    apps
// @Produce json
// @Id      listAppCronJob
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} cronjobdto.ListCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/cron-jobs [get]
func (h *AppHandler) ListAppCronJob(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeApp)
}

// GetAppCronJob Gets cron-job details
// @Summary Gets cron-job details
// @Description Gets cron-job details
// @Tags    apps
// @Produce json
// @Id      getAppCronJob
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} cronjobdto.GetCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/cron-jobs/{id} [get]
func (h *AppHandler) GetAppCronJob(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeApp)
}

// CreateAppCronJob Creates a new cron-job
// @Summary Creates a new cron-job
// @Description Creates a new cron-job
// @Tags    apps
// @Produce json
// @Id      createAppCronJob
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   body body cronjobdto.CreateCronJobReq true "request data"
// @Success 201 {object} cronjobdto.CreateCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/cron-jobs [post]
func (h *AppHandler) CreateAppCronJob(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeApp)
}

// UpdateAppCronJob Updates a cron-job
// @Summary Updates a cron-job
// @Description Updates a cron-job
// @Tags    apps
// @Produce json
// @Id      updateAppCronJob
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Param   body body cronjobdto.UpdateCronJobReq true "request data"
// @Success 200 {object} cronjobdto.UpdateCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/cron-jobs/{id} [put]
func (h *AppHandler) UpdateAppCronJob(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeApp)
}

// UpdateAppCronJobMeta Updates cron-job meta
// @Summary Updates cron-job meta
// @Description Updates cron-job meta
// @Tags    apps
// @Produce json
// @Id      updateAppCronJobMeta
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Param   body body cronjobdto.UpdateCronJobMetaReq true "request data"
// @Success 200 {object} cronjobdto.UpdateCronJobMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/cron-jobs/{id}/meta [put]
func (h *AppHandler) UpdateAppCronJobMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeCronJob, base.SettingScopeApp)
}

// DeleteAppCronJob Deletes cron-job
// @Summary Deletes cron-job
// @Description Deletes cron-job
// @Tags    apps
// @Produce json
// @Id      deleteAppCronJob
// @Param   projectID path string true "project ID"
// @Param   appID path string true "app ID"
// @Param   id path string true "provider ID"
// @Success 200 {object} cronjobdto.DeleteCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /projects/{projectID}/apps/{appID}/cron-jobs/{id} [delete]
func (h *AppHandler) DeleteAppCronJob(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeApp)
}

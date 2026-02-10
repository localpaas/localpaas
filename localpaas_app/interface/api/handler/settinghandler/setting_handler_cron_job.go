package settinghandler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	_ "github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc/cronjobdto"
)

// ListCronJob Lists cron-jobs
// @Summary Lists cron-jobs
// @Description Lists cron-jobs
// @Tags    settings
// @Produce json
// @Id      listSettingCronJob
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} cronjobdto.ListCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cron-jobs [get]
func (h *SettingHandler) ListCronJob(ctx *gin.Context) {
	h.ListSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeGlobal)
}

// GetCronJob Gets cron-job details
// @Summary Gets cron-job details
// @Description Gets cron-job details
// @Tags    settings
// @Produce json
// @Id      getSettingCronJob
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cronjobdto.GetCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cron-jobs/{itemID} [get]
func (h *SettingHandler) GetCronJob(ctx *gin.Context) {
	h.GetSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeGlobal)
}

// CreateCronJob Creates a new cron-job
// @Summary Creates a new cron-job
// @Description Creates a new cron-job
// @Tags    settings
// @Produce json
// @Id      createSettingCronJob
// @Param   body body cronjobdto.CreateCronJobReq true "request data"
// @Success 201 {object} cronjobdto.CreateCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cron-jobs [post]
func (h *SettingHandler) CreateCronJob(ctx *gin.Context) {
	h.CreateSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeGlobal)
}

// UpdateCronJob Updates cron-job
// @Summary Updates cron-job
// @Description Updates cron-job
// @Tags    settings
// @Produce json
// @Id      updateSettingCronJob
// @Param   itemID path string true "setting ID"
// @Param   body body cronjobdto.UpdateCronJobReq true "request data"
// @Success 200 {object} cronjobdto.UpdateCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cron-jobs/{itemID} [put]
func (h *SettingHandler) UpdateCronJob(ctx *gin.Context) {
	h.UpdateSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeGlobal)
}

// UpdateCronJobMeta Updates cron-job meta
// @Summary Updates cron-job meta
// @Description Updates cron-job meta
// @Tags    settings
// @Produce json
// @Id      updateSettingCronJobMeta
// @Param   itemID path string true "setting ID"
// @Param   body body cronjobdto.UpdateCronJobMetaReq true "request data"
// @Success 200 {object} cronjobdto.UpdateCronJobMetaResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cron-jobs/{itemID}/meta [put]
func (h *SettingHandler) UpdateCronJobMeta(ctx *gin.Context) {
	h.UpdateSettingMeta(ctx, base.ResourceTypeCronJob, base.SettingScopeGlobal)
}

// DeleteCronJob Deletes cron-job
// @Summary Deletes cron-job
// @Description Deletes cron-job
// @Tags    settings
// @Produce json
// @Id      deleteSettingCronJob
// @Param   itemID path string true "setting ID"
// @Success 200 {object} cronjobdto.DeleteCronJobResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /settings/cron-jobs/{itemID} [delete]
func (h *SettingHandler) DeleteCronJob(ctx *gin.Context) {
	h.DeleteSetting(ctx, base.ResourceTypeCronJob, base.SettingScopeGlobal)
}

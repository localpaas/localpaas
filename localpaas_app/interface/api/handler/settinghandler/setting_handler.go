package settinghandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/cronjobuc"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/secretuc"
)

type SettingHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	secretUC    *secretuc.SecretUC
	cronJobUC   *cronjobuc.CronJobUC
}

func NewSettingHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	secretUC *secretuc.SecretUC,
	cronJobUC *cronjobuc.CronJobUC,
) *SettingHandler {
	return &SettingHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		secretUC:    secretUC,
		cronJobUC:   cronJobUC,
	}
}

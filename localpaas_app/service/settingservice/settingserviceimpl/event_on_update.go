package settingserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

func (s *service) OnUpdate(
	ctx context.Context,
	db database.IDB,
	event *settingservice.UpdateEvent,
) (err error) {
	// Remove healthcheck cache if the update may relate
	if event.Setting.IsTypeIn(base.SettingTypeHealthcheck, base.SettingTypeIMService, base.SettingTypeEmail) {
		err = s.healthcheckSettingsRepo.Del(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	// Save SSL cert/key files in a directory for using later
	if event.Setting.Type == base.SettingTypeSSLCert {
		err = s.sslService.WriteCertFiles(true, event.Setting)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

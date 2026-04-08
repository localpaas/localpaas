package settingserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

func (s *service) OnCreate(
	ctx context.Context,
	db database.IDB,
	event *settingservice.CreateEvent,
) (err error) {
	// Save SSL cert/key files in a directory for using later
	if event.Setting.Type == base.SettingTypeSSLCert {
		err = s.PersistSSLCertFiles(true, event.Setting)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

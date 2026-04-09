package settingserviceimpl

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
)

const (
	sslCertSettingName        = "SSL certificate settings"
	sslCertTypeDefault        = base.SSLCertTypeLetsEncrypt
	sslCertKeyTypeDefault     = base.SSLKeyTypeECP256
	sslCertValidPeriodDefault = time.Hour * 24 * 365 // For self-signed certs only
	sslCertAutoRenewDefault   = true
)

func (s *service) initDefaultSSLCertSettings(
	ctx context.Context,
	db database.IDB,
	timeNow time.Time,
) (err error) {
	sslCertSetting := &entity.Setting{
		ID:              gofn.Must(ulid.NewStringULID()),
		Scope:           base.SettingScopeGlobal,
		Type:            base.SettingTypeSSLCertSettings,
		Status:          base.SettingStatusActive,
		Name:            sslCertSettingName,
		AvailInProjects: true,
		Default:         true,
		Version:         entity.CurrentSSLCertSettingsVersion,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
	}
	sslCert := &entity.SSLCertSettings{
		CertType:    sslCertTypeDefault,
		KeyType:     sslCertKeyTypeDefault,
		ValidPeriod: timeutil.Duration(sslCertValidPeriodDefault),
		RootDomain:  config.Current.RootDomain,
		Email:       config.Current.AdminAccount.Email,
		AutoRenew:   sslCertAutoRenewDefault,
	}
	sslCertSetting.MustSetData(sslCert)

	err = s.settingRepo.Insert(ctx, db, sslCertSetting)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

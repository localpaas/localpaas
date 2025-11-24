package appuc

import (
	"context"
	"path/filepath"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/fileutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
	"github.com/localpaas/localpaas/localpaas_app/usecase/appuc/appdto"
	"github.com/localpaas/localpaas/services/letsencrypt"
)

func (uc *AppUC) ObtainDomainSsl(
	ctx context.Context,
	auth *basedto.Auth,
	req *appdto.ObtainDomainSslReq,
) (*appdto.ObtainDomainSslResp, error) {
	email := gofn.Coalesce(req.Email, config.Current.SSL.LeUserEmail)
	leClient, err := letsencrypt.NewClient(email, req.KeySize, config.Current.DataPathNginxShareDomains())
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	certificates, err := leClient.ObtainCertificate(ctx, []string{req.Domain})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	appData := &obtainSslData{
		ObtainedCerts: certificates,
		Email:         email,
	}

	err = transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		err = uc.loadAppDataForObtainDomainSsl(ctx, uc.db, req, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingAppData{}
		uc.preparePersistingDomainSslData(req, appData, persistingData)

		err = uc.appService.PersistAppData(ctx, db, &persistingData.PersistingAppData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		err = uc.applyDomainSsl(ctx, appData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &appdto.ObtainDomainSslResp{
		Data: &basedto.ObjectIDResp{ID: appData.SslCert.Setting.ID},
	}, nil
}

type obtainSslData struct {
	App           *entity.App
	HttpSettings  *entity.AppHttpSettings
	SslCert       *entity.Ssl
	ObtainedCerts *certificate.Resource
	Email         string
}

func (uc *AppUC) loadAppDataForObtainDomainSsl(
	ctx context.Context,
	db database.IDB,
	req *appdto.ObtainDomainSslReq,
	data *obtainSslData,
) error {
	app, err := uc.appRepo.GetByID(ctx, db, req.ProjectID, req.AppID,
		bunex.SelectFor("UPDATE"),
		bunex.SelectRelation("Settings",
			bunex.SelectWhere("setting.type = ?", base.SettingTypeAppHttp),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if app.Status != base.AppStatusActive {
		return apperrors.Wrap(apperrors.ErrResourceInactive)
	}
	data.App = app

	dbSetting := app.GetSettingByType(base.SettingTypeAppHttp)
	if dbSetting == nil {
		return apperrors.NewNotFound("AppHttpSetting")
	}

	httpSettings, err := dbSetting.ParseAppHttpSettings()
	if err != nil {
		return apperrors.Wrap(err)
	}

	domainSettings := httpSettings.GetDomain(req.Domain)
	if domainSettings != nil && domainSettings.SslCert.ID != "" {
		return apperrors.New(apperrors.ErrAlreadyExist).
			WithMsgLog("ssl for domain '%s' already exists", req.Domain)
	}

	data.HttpSettings = httpSettings

	return nil
}

func (uc *AppUC) preparePersistingDomainSslData(
	req *appdto.ObtainDomainSslReq,
	data *obtainSslData,
	persistingData *persistingAppData,
) {
	timeNow := timeutil.NowUTC()
	dbSsl := &entity.Setting{
		ID:        gofn.Must(ulid.NewStringULID()),
		Type:      base.SettingTypeSsl,
		Status:    base.SettingStatusActive,
		Name:      req.Domain,
		Kind:      string(base.SslProviderLetsEncrypt),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	ssl := &entity.Ssl{
		Setting:     dbSsl,
		Certificate: string(data.ObtainedCerts.Certificate),
		PrivateKey:  string(data.ObtainedCerts.PrivateKey),
		KeySize:     req.KeySize,
		Provider:    base.SslProviderLetsEncrypt,
		Email:       data.Email,
	}
	data.SslCert = ssl

	dbSsl.MustSetData(ssl.MustEncrypt())
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, dbSsl)

	httpSettings := data.HttpSettings
	domainSettings := httpSettings.GetDomain(req.Domain)
	if domainSettings == nil {
		domainSettings = &entity.AppDomain{Domain: req.Domain}
		httpSettings.Domains = append(httpSettings.Domains, domainSettings)
	}
	domainSettings.SslCert.ID = dbSsl.ID

	// Enables the HTTP settings
	httpSettings.Enabled = true
	httpSettings.Setting.MustSetData(httpSettings)

	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, httpSettings.Setting)
}

func (uc *AppUC) applyDomainSsl(
	ctx context.Context,
	data *obtainSslData,
) error {
	saveDir := filepath.Join(config.Current.DataPathCerts(), data.SslCert.Setting.ID)
	err := fileutil.WriteCerts(data.ObtainedCerts.Certificate, data.ObtainedCerts.PrivateKey, saveDir,
		"certificate.crt", "private.key")
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = uc.nginxService.ApplyAppConfig(ctx, data.App, data.HttpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = uc.networkService.UpdateAppGlobalRoutingNetwork(ctx, data.App, data.HttpSettings)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

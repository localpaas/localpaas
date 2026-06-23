package sslserviceimpl

import (
	"github.com/go-acme/lego/v5/lego"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/ssl/acme"
)

func (s *service) GetAcmeClient(
	sslSetting *entity.Setting,
	refObjects *entity.RefObjects,
) (_ *acme.Client, err error) {
	sslCert := sslSetting.MustAsSSLCert()
	acmeCfg := &acme.ACMEConfig{
		Email: sslCert.Email,
	}

	if sslCert.Provider.ID != "" {
		providerSetting := refObjects.RefSettings[sslCert.Provider.ID]
		if providerSetting == nil {
			return nil, apperrors.NewNotFound(apperrors.Fmt("SSL provider '%v'", sslCert.Provider.ID))
		}
		provider := providerSetting.MustAsSSLProvider()

		switch sslCert.CertType {
		case base.SSLCertTypeLetsEncrypt:
			// Do nothing for now
		case base.SSLCertTypeZeroSSL:
			acmeCfg.CACode = lego.CodeZeroSSL
			acmeCfg.EABKid = provider.ZeroSSL.EABKid
			acmeCfg.EABHmacKey, err = provider.ZeroSSL.EABHmacKey.GetPlain()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
		case base.SSLCertTypeGoogleTrust:
			acmeCfg.CACode = lego.CodeGoogleTrust
			acmeCfg.EABKid = provider.GoogleTrust.EABKid
			acmeCfg.EABHmacKey, err = provider.GoogleTrust.EABHmacKey.GetPlain()
			if err != nil {
				return nil, apperrors.Wrap(err)
			}
		case base.SSLCertTypeSelfSigned, base.SSLCertTypeCustom:
			// Do nothing
		}
	}

	if sslCert.AcmeProvider.ID != "" {
		acmeProviderSetting := refObjects.RefSettings[sslCert.AcmeProvider.ID]
		if acmeProviderSetting == nil {
			return nil, apperrors.NewNotFound(apperrors.Fmt("ACME provider '%v'", sslCert.AcmeProvider.ID))
		}
		// NOTE: now there is only DNS-01 provider type of the setting
		acmeDnsProvider := acmeProviderSetting.MustAsAcmeDnsProvider()

		acmeCfg.DNS01Provider, err = acme.NewDNS01Provider(base.AcmeDnsProvider(acmeProviderSetting.Kind),
			acmeDnsProvider)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	// If no DNS-01 provider is set, fallback to HTTP-01
	if acmeCfg.DNS01Provider == nil {
		acmeCfg.HTTP01Provider, err = acme.NewHTTP01Provider(config.Current.DataPathSslAcme().AbsPath())
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	acmeClient, err := acme.NewClient(acmeCfg)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return acmeClient, nil
}

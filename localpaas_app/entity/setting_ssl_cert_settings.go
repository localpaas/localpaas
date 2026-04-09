package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	CurrentSSLCertSettingsVersion = 1
)

var _ = registerSettingParser(base.SettingTypeSSLCertSettings, &sslCertSettingsParser{})

type sslCertSettingsParser struct {
}

func (s *sslCertSettingsParser) New() SettingData {
	return &SSLCertSettings{}
}

type SSLCertSettings struct {
	CertType    base.SSLCertType  `json:"certType"`
	KeyType     base.SSLKeyType   `json:"keyType"`
	ValidPeriod timeutil.Duration `json:"validPeriod,omitempty"`
	RootDomain  string            `json:"rootDomain,omitempty"`
	Email       string            `json:"email"`
	AutoRenew   bool              `json:"autoRenew,omitempty"`
}

func (s *SSLCertSettings) GetType() base.SettingType {
	return base.SettingTypeSSLCertSettings
}

func (s *SSLCertSettings) GetRefObjectIDs() *RefObjectIDs {
	refIDs := &RefObjectIDs{}
	return refIDs
}

func (s *SSLCertSettings) Migrate(setting *Setting) (hasChange bool, err error) {
	if setting.Version == CurrentSSLCertSettingsVersion {
		return false, nil
	}
	if setting.Version > CurrentSSLCertSettingsVersion {
		return false, apperrors.New(apperrors.ErrDataVerNewerThanSystemVer)
	}

	// TODO: add migration if we make any change

	setting.Version = CurrentSSLCertSettingsVersion
	setting.UpdateVer++
	setting.MustSetData(s)
	return true, nil
}

func (s *Setting) AsSSLCertSettings() (*SSLCertSettings, error) {
	return parseSettingAs[*SSLCertSettings](s)
}

func (s *Setting) MustAsSSLCertSettings() *SSLCertSettings {
	return gofn.Must(s.AsSSLCertSettings())
}

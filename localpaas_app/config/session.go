package config

import "time"

const (
	AuthTypeSAML = "saml"
)

type Session struct {
	// NOTE: using `type AuthType string` causes `configor` fail to load config
	AuthType               string        `yaml:"auth_type" env:"LP_SESSION_AUTH_TYPE"`
	LastAccessUpdatePeriod time.Duration `yaml:"last_access_update_period" env:"LP_SESSION_LAST_ACCESS_UPDATE_PERIOD" default:"1m"` //nolint:lll

	SAML struct {
		IDPMetadataURL    string   `yaml:"idp_metadata_url" env:"LP_SESSION_SAML_IDP_METADATA_URL"`
		SPIssuer          string   `yaml:"sp_issuer" env:"LP_SESSION_SAML_SP_ISSUER"`
		SPAudienceURI     string   `yaml:"sp_audience_uri" env:"LP_SESSION_SAML_SP_AUDIENCE_URI"`
		SPPrivateKeyFile  string   `yaml:"sp_private_key_file" env:"LP_SESSION_SAML_SP_PRIVATE_KEY_FILE"`
		SPCertificateFile string   `yaml:"sp_certificate_file" env:"LP_SESSION_SAML_SP_CERTIFICATE_FILE"`
		AttrMappings      []string `yaml:"attr_mappings" env:"LP_SESSION_SAML_ATTR_MAPPINGS"`
	} `yaml:"saml"`

	BasicAuth struct {
		Username string `yaml:"username" env:"LP_SESSION_BASIC_AUTH_USERNAME"`
		Password string `yaml:"password" env:"LP_SESSION_BASIC_AUTH_PASSWORD"`
	} `yaml:"basic_auth"`

	JWT struct {
		Secret          string        `yaml:"secret" env:"LP_SESSION_JWT_SECRET"`
		AccessTokenExp  time.Duration `yaml:"access_token_exp" env:"LP_SESSION_JWT_ACCESS_TOKEN_EXP" default:"12h"`
		RefreshTokenExp time.Duration `yaml:"refresh_token_exp" env:"LP_SESSION_JWT_REFRESH_TOKEN_EXP" default:"24h"`
	} `yaml:"jwt"`

	MFA struct {
		PasscodeTimeout     time.Duration `yaml:"passcode_timeout" env:"LP_SESSION_MFA_PASSCODE_TIMEOUT" default:"10m"`
		PasscodeResendAfter time.Duration `yaml:"passcode_resend_after" env:"LP_SESSION_MFA_PASSCODE_RESEND_AFTER" default:"60s"` //nolint:lll
		DeviceTrustedPeriod time.Duration `yaml:"device_trusted_period" env:"LP_SESSION_MFA_DEVICE_TRUSTED_PERIOD" default:"46h"` //nolint:lll
	} `yaml:"mfa"`
}

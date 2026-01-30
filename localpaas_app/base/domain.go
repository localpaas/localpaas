package base

type SSLProvider string

const (
	SSLProviderLetsEncrypt SSLProvider = "letsencrypt"
	SSLProviderCustom      SSLProvider = "custom"
)

var (
	AllSSLProviders = []SSLProvider{SSLProviderLetsEncrypt, SSLProviderCustom}
)

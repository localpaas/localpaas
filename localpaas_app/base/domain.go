package base

type SslProvider string

const (
	SslProviderLetsEncrypt SslProvider = "letsencrypt"
	SslProviderCustom      SslProvider = "custom"
)

var (
	AllSslProviders = []SslProvider{SslProviderLetsEncrypt, SslProviderCustom}
)

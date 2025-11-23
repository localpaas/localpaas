package base

type SslProvider string

const (
	SslProviderLetsEncrypt SslProvider = "letsencrypt"
)

var (
	AllSslProviders = []SslProvider{SslProviderLetsEncrypt}
)

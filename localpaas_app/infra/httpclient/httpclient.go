package httpclient

import "net/http"

var (
	DefaultClient = func() *http.Client {
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
	}()
)

// TODO: supports socks5 proxy

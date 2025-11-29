package slack

import (
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/infra/httpclient"
)

type Client struct {
	httpClient *http.Client
}

func (c *Client) getHttpClient() *http.Client {
	if c.httpClient == nil {
		return httpclient.DefaultClient
	}
	return c.httpClient
}

func NewClient() *Client {
	client := &Client{}
	return client
}

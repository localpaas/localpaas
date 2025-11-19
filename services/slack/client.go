package slack

import (
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	client := &Client{
		httpClient: http.DefaultClient,
	}
	return client
}

package discord

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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

func parseWebhookURL(webhook string) (webhookID, token string, err error) {
	dcUrl, err := url.Parse(webhook)
	if err != nil {
		return "", "", apperrors.Wrap(err)
	}
	parts := strings.Split(dcUrl.Path, "/")
	if len(parts) > 2 { //nolint:mnd
		return parts[len(parts)-2], parts[len(parts)-1], nil
	}
	return "", "", apperrors.New(apperrors.ErrParamInvalid).
		WithMsgLog("unabled to parse webhook URL")
}

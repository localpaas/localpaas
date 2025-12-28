package netutil

import (
	"net/url"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

func URLAddAuth(rawURL, user, password string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	u.User = url.UserPassword(user, password)
	return u.String(), nil
}

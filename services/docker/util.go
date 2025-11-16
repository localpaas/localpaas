package docker

import (
	"github.com/docker/docker/api/types/registry"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

func GenerateAuthHeader(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", nil
	}
	h, err := registry.EncodeAuthConfig(registry.AuthConfig{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return h, nil
}

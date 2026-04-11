package docker

import (
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
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
		return "", apperrors.NewInfra(err)
	}
	return h, nil
}

func FilterAdd(f *filters.Args, key, value string) {
	if f == nil {
		return
	}
	if f.Len() == 0 {
		*f = filters.NewArgs()
	}
	f.Add(key, value)
}

func CallRetry(
	fn func() error,
	maxRetries int,
	retryInterval time.Duration,
) error {
	retry := 0
	for {
		err := fn()
		if err == nil {
			return nil
		}
		if retry >= maxRetries {
			return apperrors.NewInfra(err)
		}
		if retryInterval > 0 {
			time.Sleep(retryInterval)
		}
		retry++
	}
}

func CallRetry2[T any](
	fn func() (T, error),
	maxRetries int,
	retryInterval time.Duration,
) (T, error) {
	retry := 0
	for {
		v, err := fn()
		if err == nil {
			return v, nil
		}
		if retry >= maxRetries {
			return v, apperrors.NewInfra(err)
		}
		if retryInterval > 0 {
			time.Sleep(retryInterval)
		}
		retry++
	}
}

func CallRetry3[T any, U any](
	fn func() (T, U, error),
	maxRetries int,
	retryInterval time.Duration,
) (T, U, error) {
	retry := 0
	for {
		v1, v2, err := fn()
		if err == nil {
			return v1, v2, nil
		}
		if retry >= maxRetries {
			return v1, v2, apperrors.NewInfra(err)
		}
		if retryInterval > 0 {
			time.Sleep(retryInterval)
		}
		retry++
	}
}

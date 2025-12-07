package github

import (
	"errors"
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

var (
	ErrGithubAppClientRequired     = errors.New("github app client required")
	ErrGithubAccessProviderInvalid = errors.New("github access provider invalid")
)

func init() {
	apperrors.RegisterStatusMapping(ErrGithubAppClientRequired, http.StatusForbidden)
}

package github

import (
	"errors"
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

var (
	ErrGithubAppClientRequired   = errors.New("github app client required")
	ErrGithubTokenClientRequired = errors.New("github token client required")
	ErrAccessProviderInvalid     = errors.New("access provider invalid")
	ErrAPICallFailed             = errors.New("api call failed")
)

func init() {
	apperrors.RegisterStatusMapping(ErrGithubAppClientRequired, http.StatusMethodNotAllowed)
	apperrors.RegisterStatusMapping(ErrGithubTokenClientRequired, http.StatusMethodNotAllowed)
	apperrors.RegisterStatusMapping(ErrAccessProviderInvalid, http.StatusNotAcceptable)
	apperrors.RegisterStatusMapping(ErrAPICallFailed, http.StatusUnprocessableEntity)
}

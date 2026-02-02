package gitea

import (
	"errors"
	"net/http"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

var (
	ErrAccessProviderInvalid = errors.New("access provider invalid")
)

func init() {
	apperrors.RegisterStatusMapping(ErrAccessProviderInvalid, http.StatusNotAcceptable)
}

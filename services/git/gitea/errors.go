package gitea

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

var (
	ErrAccessProviderInvalid = apperrors.NewErr(apperrors.ErrArgumentInvalid, "access provider invalid")
)

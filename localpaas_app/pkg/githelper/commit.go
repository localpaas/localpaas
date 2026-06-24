package githelper

import (
	"github.com/localpaas/localpaas/localpaas_app/pkg/githelper/validation"
)

func IsCommitHash(hash string) bool {
	return validation.IsCommitHash(hash)
}

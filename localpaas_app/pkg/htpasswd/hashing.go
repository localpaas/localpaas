package htpasswd

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/localpaas/localpaas/localpaas_app/pkg/tracerr"
)

func hashBcrypt(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return string(passwordBytes), nil
}

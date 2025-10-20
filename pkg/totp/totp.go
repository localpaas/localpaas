package totp

import (
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/localpaas/localpaas/pkg/tracerr"
)

func GenerateSecret() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{})
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return key.Secret(), nil
}

func GeneratePassCode(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return code, nil
}

func VerifyCode(passCode, secret string) bool {
	return totp.Validate(passCode, secret)
}

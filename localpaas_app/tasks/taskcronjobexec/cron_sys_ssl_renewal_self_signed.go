package taskcronjobexec

import (
	"context"
	"crypto/x509/pkix"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	sslSelfSignedRenewablePeriod = time.Hour * 24 * 30 // 30 days
)

func (e *Executor) sslRenewSelfSignedCert(
	_ context.Context,
	ssl *entity.SSLCert,
	_ *sslRenewalTaskData,
) (err error) {
	if !ssl.AutoRenew {
		return nil
	}

	notBefore := timeutil.NowUTC()
	notAfter := notBefore.Add(ssl.ValidPeriod.ToDuration())

	certBytes, keyBytes, err := e.sslService.GenerateCert(&pkix.Name{CommonName: ssl.Domain}, ssl.KeyType,
		notBefore, notAfter, false)
	if err != nil {
		return apperrors.Wrap(err)
	}

	ssl.Certificate = reflectutil.UnsafeBytesToStr(certBytes)
	ssl.PrivateKey = entity.NewEncryptedField(reflectutil.UnsafeBytesToStr(keyBytes))
	ssl.ExpireAt = notAfter
	ssl.RenewableFrom = ssl.ExpireAt.Add(-sslSelfSignedRenewablePeriod)
	ssl.NotifyFrom = ssl.RenewableFrom

	return nil
}

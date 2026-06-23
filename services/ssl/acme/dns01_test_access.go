package acme

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func DNS01ProviderTestAccess(
	ctx context.Context,
	providerKind base.AcmeDnsProvider,
	dnsConfig *entity.AcmeDnsProvider,
	testDomain string,
) (err error) {
	provider, err := NewDNS01Provider(providerKind, dnsConfig)
	if err != nil {
		return apperrors.New(err)
	}
	err = provider.Present(ctx, testDomain, "test", "test")
	if err != nil {
		return apperrors.New(err)
	}
	err = provider.CleanUp(ctx, testDomain, "test", "test")
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

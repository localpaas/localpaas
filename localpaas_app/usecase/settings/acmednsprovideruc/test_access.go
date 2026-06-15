package acmednsprovideruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/acmednsprovideruc/acmednsproviderdto"
	"github.com/localpaas/localpaas/services/ssl/acme"
)

func (uc *UC) TestProviderAccess(
	ctx context.Context,
	auth *basedto.Auth,
	req *acmednsproviderdto.TestProviderAccessReq,
) (*acmednsproviderdto.TestProviderAccessResp, error) {
	err := acme.DNS01ProviderTestAccess(ctx, req.Kind, req.ToEntity(), req.TestDomain)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return &acmednsproviderdto.TestProviderAccessResp{}, nil
}

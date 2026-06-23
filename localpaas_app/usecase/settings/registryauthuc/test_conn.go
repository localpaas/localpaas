package registryauthuc

import (
	"context"

	"github.com/moby/moby/client"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/registryauthuc/registryauthdto"
)

func (uc *UC) TestRegistryAuthConn(
	ctx context.Context,
	auth *basedto.Auth,
	req *registryauthdto.TestRegistryAuthConnReq,
) (*registryauthdto.TestRegistryAuthConnResp, error) {
	_, err := uc.dockerManager.RegistryLogin(ctx, func(opts *client.RegistryLoginOptions) {
		opts.Username = req.Username
		opts.Password = req.Password
		opts.ServerAddress = req.Address
	})
	if err != nil {
		return nil, apperrors.New(err)
	}

	return &registryauthdto.TestRegistryAuthConnResp{}, nil
}

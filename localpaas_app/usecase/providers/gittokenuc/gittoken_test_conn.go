package gittokenuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/gittokenuc/gittokendto"
)

func (uc *GitTokenUC) TestGitTokenConn(
	ctx context.Context,
	auth *basedto.Auth,
	req *gittokendto.TestGitTokenConnReq,
) (*gittokendto.TestGitTokenConnResp, error) {
	// TODO: add implementation
	return &gittokendto.TestGitTokenConnResp{}, nil
}

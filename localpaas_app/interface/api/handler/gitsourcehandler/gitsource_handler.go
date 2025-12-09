package gitsourcehandler

import (
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler"
	"github.com/localpaas/localpaas/localpaas_app/interface/api/handler/authhandler"
	"github.com/localpaas/localpaas/localpaas_app/usecase/gitsourceuc"
)

type GitSourceHandler struct {
	*handler.BaseHandler
	authHandler *authhandler.AuthHandler
	gitSourceUC *gitsourceuc.GitSourceUC
}

func NewGitSourceHandler(
	baseHandler *handler.BaseHandler,
	authHandler *authhandler.AuthHandler,
	gitSourceUC *gitsourceuc.GitSourceUC,
) *GitSourceHandler {
	return &GitSourceHandler{
		BaseHandler: baseHandler,
		authHandler: authHandler,
		gitSourceUC: gitSourceUC,
	}
}

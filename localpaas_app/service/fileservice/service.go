package fileservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity/appentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type Service interface {
	GetDownloadURL(ctx context.Context, db database.IDB, auth *basedto.Auth, req *GetDownloadURLReq) (
		*GetDownloadURLResp, error)
	ParseFileDownloadToken(token string) (*appentity.FileDownloadTokenClaims, error)
}

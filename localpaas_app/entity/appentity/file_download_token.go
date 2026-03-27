package appentity

import (
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

type FileDownloadTokenClaims struct {
	jwtsession.BaseClaims
	FileId       string `json:"fileId"`
	UserId       string `json:"userId"`
	RequireLogin bool   `json:"requireLogin,omitempty"`
}

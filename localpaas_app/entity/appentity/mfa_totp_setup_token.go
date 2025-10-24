package appentity

import "github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"

type MFATotpSetupTokenClaims struct {
	jwtsession.BaseClaims
	UserID string `json:"userId"`
	Secret string `json:"secret"`
}

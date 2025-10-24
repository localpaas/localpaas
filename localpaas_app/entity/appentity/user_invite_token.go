package appentity

import (
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

type UserInviteTokenClaims struct {
	jwtsession.BaseClaims
	UserID string `json:"userId"`
}

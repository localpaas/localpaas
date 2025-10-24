package appentity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

type MFATokenClaims struct {
	jwtsession.BaseClaims
	UserID          string       `json:"userId"`
	MFAType         base.MFAType `json:"mfaType"`
	TrustedDeviceID string       `json:"deviceId,omitempty"`
}

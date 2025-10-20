package sessiondto

import (
	"time"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

type BaseCreateSessionReq struct {
	User *entity.User
}

type BaseCreateSessionResp struct {
	AccessToken     string    `json:"accessToken"`
	AccessTokenExp  time.Time `json:"accessTokenExp"`
	RefreshToken    string    `json:"-"`
	RefreshTokenExp time.Time `json:"-"`
}

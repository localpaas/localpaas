package basedto

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jwtsession"
)

type User struct {
	*entity.User
	AuthClaims *jwtsession.AuthClaims
}

func (u *User) Entity() *entity.User {
	if u != nil {
		return u.User
	}
	return nil
}

type Auth struct {
	User *User
}

package emailservice

import "github.com/localpaas/localpaas/localpaas_app/entity"

type EmailDataPasswordReset struct {
	Email             *entity.Email
	Recipients        []string
	Subject           string
	ResetPasswordLink string
}

type EmailDataUserInvite struct {
	Email          *entity.Email
	Recipients     []string
	Subject        string
	InviterName    string
	UserSignupLink string
}

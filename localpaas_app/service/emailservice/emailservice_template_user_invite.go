package emailservice

import (
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/services/email"
)

type EmailDataUserInvite struct {
	Email          *entity.Email
	Recipients     []string
	Subject        string
	InviterName    string
	UserSignupLink string
}

func (s *emailService) SendMailUserInvite(
	ctx context.Context,
	db database.IDB,
	data *EmailDataUserInvite,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeUserInvite)
	if err != nil {
		return apperrors.Wrap(err)
	}

	subject := gofn.Coalesce(data.Subject, "You’ve been invited to join LocalPaaS")

	content := template.ExecuteString(map[string]any{
		"inviter_name":     data.InviterName,
		"user_signup_link": data.UserSignupLink,
	})

	err = email.SendMail(ctx, data.Email, data.Recipients, subject, content)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

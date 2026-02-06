package emailservice

import (
	"bytes"
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

	buf := bytes.NewBuffer(make([]byte, 0, buffSizeMd))
	err = template.Execute(buf, *data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	subject := gofn.Coalesce(data.Subject, "You’ve been invited to join LocalPaaS")
	err = email.SendMail(ctx, data.Email, data.Recipients, subject, buf.String())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

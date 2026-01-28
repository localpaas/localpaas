package email

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/email/http"
	"github.com/localpaas/localpaas/services/email/smtp"
)

func SendMail(
	ctx context.Context,
	email *entity.Email,
	recipients []string,
	subject string,
	content string,
) (err error) {
	switch { //nolint
	case email.SMTP != nil:
		err = smtp.SendMail(ctx, email.SMTP, recipients, subject, content)
	case email.HTTP != nil:
		err = http.SendMail(ctx, email.HTTP, recipients, subject, content)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

package emailuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
	"github.com/localpaas/localpaas/services/email/http"
	"github.com/localpaas/localpaas/services/email/smtp"
)

func (uc *EmailUC) TestSendEmail(
	ctx context.Context,
	auth *basedto.Auth,
	req *emaildto.TestSendEmailReq,
) (_ *emaildto.TestSendEmailResp, err error) {
	email := req.ToEntity()
	switch req.Kind {
	case base.EmailKindSMTP:
		err = smtp.SendMail(ctx, email.SMTP, []string{req.TestRecipient}, req.TestSubject, req.TestContent)
	case base.EmailKindHTTP:
		err = http.SendMail(ctx, email.HTTP, []string{req.TestRecipient}, req.TestSubject, req.TestContent)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &emaildto.TestSendEmailResp{}, nil
}

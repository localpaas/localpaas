package emailuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/emailuc/emaildto"
	"github.com/localpaas/localpaas/services/email/smtp"
)

func (uc *EmailUC) TestSendEmail(
	ctx context.Context,
	auth *basedto.Auth,
	req *emaildto.TestSendEmailReq,
) (*emaildto.TestSendEmailResp, error) {
	switch {
	case req.SMTP != nil:
		if err := uc.testSendSmtpEmail(ctx, req); err != nil {
			return nil, apperrors.Wrap(err)
		}
	case req.HTTP != nil:
		if err := uc.testSendHttpEmail(ctx, req); err != nil {
			return nil, apperrors.Wrap(err)
		}
	}

	return &emaildto.TestSendEmailResp{}, nil
}

func (uc *EmailUC) testSendSmtpEmail(
	ctx context.Context,
	req *emaildto.TestSendEmailReq,
) error {
	conf := &entity.SMTPConf{
		Host:        req.SMTP.Host,
		Port:        req.SMTP.Port,
		Username:    req.SMTP.Username,
		DisplayName: req.SMTP.DisplayName,
		Password:    entity.NewEncryptedField(req.SMTP.Password),
		SSL:         req.SMTP.SSL,
	}

	err := smtp.SendMail(ctx, conf, []string{req.TestRecipient}, req.TestSubject, req.TestContent)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (uc *EmailUC) testSendHttpEmail(
	_ context.Context,
	_ *emaildto.TestSendEmailReq,
) error {
	// TODO: implement this
	return nil
}

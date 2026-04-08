package notificationserviceimpl

import (
	"bytes"
	"context"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/strutil"
	"github.com/localpaas/localpaas/localpaas_app/service/notificationservice"
	"github.com/localpaas/localpaas/services/email"
)

func (s *service) EmailSendSSLRenewalNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.EmailMsgDataSSLRenewalNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeEmail, TemplateSSLRenewalNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, buffSizeMd))
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	subject := gofn.Coalesce(data.Subject, "SSL renewal notification")
	err = email.SendMail(ctx, data.Email, data.Recipients, subject, buf.String())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) SlackSendSSLRenewalNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.SlackMsgDataSSLRenewalNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeSlack, TemplateSSLRenewalNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, buffSizeXs))
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.slackSendMsg(ctx, data.Setting, strutil.RemoveEmptyLines(buf.String(), false))
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) DiscordSendSSLRenewalNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.DiscordMsgDataSSLRenewalNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeDiscord, TemplateSSLRenewalNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, buffSizeXs))
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.discordSendMsg(ctx, data.Setting, strutil.RemoveEmptyLines(buf.String(), false))
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

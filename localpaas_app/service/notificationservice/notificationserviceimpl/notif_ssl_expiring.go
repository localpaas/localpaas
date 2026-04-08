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

func (s *service) EmailSendSSLExpiringNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.EmailMsgDataSSLExpiringNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeEmail, TemplateSSLExpiringNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, buffSizeMd))
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	subject := gofn.Coalesce(data.Subject, "Your SSL is expiring")
	err = email.SendMail(ctx, data.Email, data.Recipients, subject, buf.String())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) SlackSendSSLExpiringNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.SlackMsgDataSSLExpiringNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeSlack, TemplateSSLExpiringNotification)
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

func (s *service) DiscordSendSSLExpiringNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.DiscordMsgDataSSLExpiringNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeDiscord, TemplateSSLExpiringNotification)
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

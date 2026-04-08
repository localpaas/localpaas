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

const (
	buffSizeXs = 1000
	buffSizeMd = 5000
)

func (s *service) EmailSendAppDeploymentNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.EmailMsgDataAppDeploymentNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeEmail, TemplateAppDeploymentNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, buffSizeMd))
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	subject := gofn.Coalesce(data.Subject, "Deployment notification")
	err = email.SendMail(ctx, data.Email, data.Recipients, subject, buf.String())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *service) SlackSendAppDeploymentNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.SlackMsgDataAppDeploymentNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeSlack, TemplateAppDeploymentNotification)
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

func (s *service) DiscordSendAppDeploymentNotification(
	ctx context.Context,
	db database.IDB,
	data *notificationservice.DiscordMsgDataAppDeploymentNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeDiscord, TemplateAppDeploymentNotification)
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

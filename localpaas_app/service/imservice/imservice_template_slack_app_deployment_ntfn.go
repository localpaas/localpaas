package imservice

import (
	"bytes"
	"context"
	"time"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type SlackMsgDataAppDeploymentNotification struct {
	Slack *entity.Slack

	ProjectName   string
	AppName       string
	Succeeded     bool
	Method        base.DeploymentMethod
	RepoURL       string
	RepoRef       string
	Image         string
	SourceArchive string
	Duration      time.Duration
	DashboardLink string
}

func (s *imService) SlackSendAppDeploymentNotification(
	ctx context.Context,
	db database.IDB,
	data *SlackMsgDataAppDeploymentNotification,
) error {
	template, err := s.SlackGetTemplate(ctx, db, SlackTemplateAppDeploymentNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 5000)) //nolint
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.slackSendMsg(ctx, data.Slack, buf.String())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

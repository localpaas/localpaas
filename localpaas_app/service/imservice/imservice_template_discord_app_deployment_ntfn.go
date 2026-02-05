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

type DiscordMsgDataAppDeploymentNotification struct {
	Discord *entity.Discord

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

func (s *imService) DiscordSendAppDeploymentNotification(
	ctx context.Context,
	db database.IDB,
	data *DiscordMsgDataAppDeploymentNotification,
) error {
	template, err := s.DiscordGetTemplate(ctx, db, DiscordTemplateAppDeploymentNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 5000)) //nolint
	err = template.Execute(buf, data)
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = s.discordSendMsg(ctx, data.Discord, buf.String())
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

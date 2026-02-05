package emailservice

import (
	"bytes"
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/services/email"
)

type EmailDataAppDeploymentNotification struct {
	Email      *entity.Email
	Recipients []string
	Subject    string

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

func (s *emailService) SendMailAppDeploymentNotification(
	ctx context.Context,
	db database.IDB,
	data *EmailDataAppDeploymentNotification,
) error {
	template, err := s.GetTemplate(ctx, db, TemplateTypeAppDeploymentNotification)
	if err != nil {
		return apperrors.Wrap(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0, 5000)) //nolint
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

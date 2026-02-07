package notificationservice

import (
	"context"
	htmltemplate "html/template"
	"io"
	"sync"
	texttemplate "text/template"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

const (
	emailTemplateDir   = "config/email/templates/" // NOTE: must end with /
	slackTemplateDir   = "config/slack/templates/"
	discordTemplateDir = "config/discord/templates/"
)

type TemplateType string

const (
	TemplateTypeEmail   TemplateType = "email"
	TemplateTypeSlack   TemplateType = "slack"
	TemplateTypeDiscord TemplateType = "discord"
)

type TemplateName string

const (
	TemplateAppDeploymentNotification TemplateName = "app-deployment-notification"
	TemplateCronTaskNotification      TemplateName = "cron-job-notification"
)

type Template interface {
	Execute(wr io.Writer, data any) error
}

var (
	templateMap = map[TemplateType]map[TemplateName]Template{}
	mu          sync.Mutex
)

func (s *notificationService) GetTemplate(
	ctx context.Context,
	db database.IDB,
	typ TemplateType,
	name TemplateName,
) (tpl Template, err error) {
	mu.Lock()
	defer mu.Unlock()

	mapTplByName, exists := templateMap[typ]
	if !exists {
		mapTplByName = make(map[TemplateName]Template, 5) //nolint:mnd
		templateMap[typ] = mapTplByName
	}

	tpl, exists = mapTplByName[name]
	if exists {
		return tpl, nil
	}

	switch typ { //nolint
	case TemplateTypeEmail:
		tpl, err = s.loadEmailTemplate(ctx, db, name)
	case TemplateTypeSlack:
		tpl, err = s.loadSlackTemplate(ctx, db, name)
	case TemplateTypeDiscord:
		tpl, err = s.loadDiscordTemplate(ctx, db, name)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	mapTplByName[name] = tpl

	return tpl, nil
}

func (s *notificationService) loadEmailTemplate(
	_ context.Context,
	_ database.IDB,
	name TemplateName,
) (tpl Template, err error) {
	switch name { //nolint
	case TemplateAppDeploymentNotification:
		tpl, err = htmltemplate.ParseFiles(emailTemplateDir + "app_deployment_notification.html")
	case TemplateCronTaskNotification:
		tpl, err = htmltemplate.ParseFiles(emailTemplateDir + "cron_task_notification.html")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return tpl, nil
}

func (s *notificationService) loadSlackTemplate(
	_ context.Context,
	_ database.IDB,
	name TemplateName,
) (tpl Template, err error) {
	switch name { //nolint
	case TemplateAppDeploymentNotification:
		tpl, err = texttemplate.ParseFiles(slackTemplateDir + "app_deployment_notification.tpl")
	case TemplateCronTaskNotification:
		tpl, err = texttemplate.ParseFiles(slackTemplateDir + "cron_task_notification.tpl")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return tpl, nil
}

func (s *notificationService) loadDiscordTemplate(
	_ context.Context,
	_ database.IDB,
	name TemplateName,
) (tpl Template, err error) {
	switch name { //nolint
	case TemplateAppDeploymentNotification:
		tpl, err = texttemplate.ParseFiles(discordTemplateDir + "app_deployment_notification.tpl")
	case TemplateCronTaskNotification:
		tpl, err = texttemplate.ParseFiles(discordTemplateDir + "cron_task_notification.tpl")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return tpl, nil
}

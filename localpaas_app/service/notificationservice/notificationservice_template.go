package notificationservice

import (
	"context"
	"html/template"
	"sync"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
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

var (
	templateMap = map[TemplateType]map[TemplateName]*template.Template{}
	mu          sync.Mutex
)

func (s *notificationService) GetTemplate(
	ctx context.Context,
	db database.IDB,
	typ TemplateType,
	name TemplateName,
) (tpl *template.Template, err error) {
	mu.Lock()
	defer mu.Unlock()

	mapTplByName, exists := templateMap[typ]
	if !exists {
		mapTplByName = make(map[TemplateName]*template.Template, 5) //nolint:mnd
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
) (tpl *template.Template, err error) {
	switch name { //nolint
	case TemplateAppDeploymentNotification:
		tpl, err = template.ParseFiles("config/email/templates/app_deployment_notification.html")
	case TemplateCronTaskNotification:
		tpl, err = template.ParseFiles("config/email/templates/cron_task_notification.html")
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
) (tpl *template.Template, err error) {
	switch name { //nolint
	case TemplateAppDeploymentNotification:
		tpl, err = template.ParseFiles("config/slack/templates/app_deployment_notification.txt")
	case TemplateCronTaskNotification:
		tpl, err = template.ParseFiles("config/slack/templates/cron_task_notification.txt")
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
) (tpl *template.Template, err error) {
	switch name { //nolint
	case TemplateAppDeploymentNotification:
		tpl, err = template.ParseFiles("config/discord/templates/app_deployment_notification.txt")
	case TemplateCronTaskNotification:
		tpl, err = template.ParseFiles("config/discord/templates/cron_task_notification.txt")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return tpl, nil
}

package imservice

import (
	"context"
	"html/template"
	"sync"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type SlackTemplateType string

const (
	SlackTemplateAppDeploymentNotification SlackTemplateType = "app-deployment-notification"
)

var (
	slackTemplateMap = map[SlackTemplateType]*template.Template{}
	slackMtx         sync.Mutex
)

func (s *imService) SlackGetTemplate(
	_ context.Context,
	_ database.IDB,
	typ SlackTemplateType,
) (tpl *template.Template, err error) {
	slackMtx.Lock()
	defer slackMtx.Unlock()

	if tpl, exists := slackTemplateMap[typ]; exists {
		return tpl, nil
	}

	switch typ { //nolint
	case SlackTemplateAppDeploymentNotification:
		tpl, err = template.ParseFiles("config/slack_templates/app_deployment_notification.txt")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	slackTemplateMap[typ] = tpl

	return tpl, nil
}

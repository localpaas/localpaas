package imservice

import (
	"context"
	"html/template"
	"sync"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type DiscordTemplateType string

const (
	DiscordTemplateAppDeploymentNotification DiscordTemplateType = "app-deployment-notification"
)

var (
	discordTemplateMap = map[DiscordTemplateType]*template.Template{}
	discordMtx         sync.Mutex
)

func (s *imService) DiscordGetTemplate(
	_ context.Context,
	_ database.IDB,
	typ DiscordTemplateType,
) (tpl *template.Template, err error) {
	discordMtx.Lock()
	defer discordMtx.Unlock()

	if tpl, exists := discordTemplateMap[typ]; exists {
		return tpl, nil
	}

	switch typ { //nolint
	case DiscordTemplateAppDeploymentNotification:
		tpl, err = template.ParseFiles("config/discord_templates/app_deployment_notification.txt")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	discordTemplateMap[typ] = tpl

	return tpl, nil
}

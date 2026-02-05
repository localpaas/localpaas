package emailservice

import (
	"context"
	"html/template"
	"sync"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type TemplateType string

const (
	TemplateTypePasswordReset             TemplateType = "password-reset"
	TemplateTypeUserInvite                TemplateType = "user-invite"
	TemplateTypeAppDeploymentNotification TemplateType = "app-deployment-notification"
)

var (
	templateMap = map[TemplateType]*template.Template{}
	mu          sync.Mutex
)

func (s *emailService) GetTemplate(
	_ context.Context,
	_ database.IDB,
	typ TemplateType,
) (tpl *template.Template, err error) {
	mu.Lock()
	defer mu.Unlock()

	if tpl, exists := templateMap[typ]; exists {
		return tpl, nil
	}

	switch typ { //nolint
	case TemplateTypePasswordReset:
		tpl, err = template.ParseFiles("config/email_templates/password_reset.html")
	case TemplateTypeUserInvite:
		tpl, err = template.ParseFiles("config/email_templates/user_invite.html")
	case TemplateTypeAppDeploymentNotification:
		tpl, err = template.ParseFiles("config/email_templates/app_deployment_notification.html")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	templateMap[typ] = tpl

	return tpl, nil
}

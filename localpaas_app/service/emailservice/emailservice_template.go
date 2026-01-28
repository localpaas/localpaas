package emailservice

import (
	"context"
	"os"
	"sync"

	"github.com/valyala/fasttemplate"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/reflectutil"
)

type TemplateType string

const (
	TemplateTypePasswordReset TemplateType = "password-reset"
	TemplateTypeUserInvite    TemplateType = "user-invite"
)

var (
	templateMap = map[TemplateType]*fasttemplate.Template{}
	mu          sync.Mutex
)

func (s *emailService) GetTemplate(
	_ context.Context,
	_ database.IDB,
	typ TemplateType,
) (_ *fasttemplate.Template, err error) {
	mu.Lock()
	defer mu.Unlock()

	if template, exists := templateMap[typ]; exists {
		return template, nil
	}

	var data []byte
	switch typ { //nolint
	case TemplateTypePasswordReset:
		data, err = os.ReadFile("config/email_templates/password_reset.html")
	case TemplateTypeUserInvite:
		data, err = os.ReadFile("config/email_templates/user_invite.html")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	template := fasttemplate.New(reflectutil.UnsafeBytesToStr(data), "{{", "}}")
	templateMap[typ] = template

	return template, nil
}

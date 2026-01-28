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
)

type templateData struct {
	TemplateStr string
	Template    *fasttemplate.Template
}

var (
	templateMap = map[TemplateType]*templateData{}
	mu          sync.Mutex
)

func (s *emailService) GetTemplate(
	_ context.Context,
	_ database.IDB,
	typ TemplateType,
) (_ *templateData, err error) {
	mu.Lock()
	defer mu.Unlock()

	if template, exists := templateMap[typ]; exists {
		return template, nil
	}

	var data []byte
	switch typ { //nolint
	case TemplateTypePasswordReset:
		data, err = os.ReadFile("config/email_templates/password_reset.html")
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	template := &templateData{
		TemplateStr: reflectutil.UnsafeBytesToStr(data),
	}
	template.Template = fasttemplate.New(template.TemplateStr, "{{", "}}")
	templateMap[typ] = template

	return template, nil
}

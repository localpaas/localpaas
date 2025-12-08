package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

type GitToken struct {
	User    string         `json:"user"`
	Token   EncryptedField `json:"token"`
	BaseURL string         `json:"baseURL"`
}

func (s *GitToken) MustDecrypt() *GitToken {
	s.Token.MustGetPlain()
	return s
}

func (s *Setting) AsGitToken() (*GitToken, error) {
	return parseSettingAs(s, base.SettingTypeGitToken, func() *GitToken { return &GitToken{} })
}

func (s *Setting) MustAsGitToken() *GitToken {
	return gofn.Must(s.AsGitToken())
}

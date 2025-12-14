package entity

import (
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/base"
)

const (
	CurrentAPIKeyVersion = 1
)

type APIKey struct {
	KeyID        string              `json:"keyId"`
	SecretKey    HashField           `json:"secretKey"`
	AccessAction *base.AccessActions `json:"accessAction,omitempty"`
}

func (s *Setting) AsAPIKey() (*APIKey, error) {
	return parseSettingAs(s, base.SettingTypeAPIKey, func() *APIKey { return &APIKey{} })
}

func (s *Setting) MustAsAPIKey() *APIKey {
	return gofn.Must(s.AsAPIKey())
}

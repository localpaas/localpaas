package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type AppHttpSettings struct {
	Enabled bool `json:"enabled"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (s *Setting) ParseAppHttpSettings() (*AppHttpSettings, error) {
	res := &AppHttpSettings{Setting: s}
	if s != nil && s.Data != "" && s.Type == base.SettingTypeAppHttp {
		return res, s.parseData(res)
	}
	return nil, nil
}

package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
)

type APIKey struct {
	SecretKey    string          `json:"secretKey"`
	Salt         string          `json:"salt"`
	AccessAction base.ActionType `json:"accessAction,omitempty"`
}

func (s *Setting) ParseAPIKey() (*APIKey, error) {
	if s != nil && s.Data != "" {
		res := &APIKey{}
		err := s.parseData(res)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return res, nil
	}
	return nil, nil
}

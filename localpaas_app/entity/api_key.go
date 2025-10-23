package entity

import "github.com/localpaas/localpaas/localpaas_app/apperrors"

type APIKey struct {
	SecretKey  string           `json:"secretKey"`
	Salt       string           `json:"salt"`
	ActingUser APIKeyActingUser `json:"actingUser"`
}

type APIKeyActingUser struct {
	ID string `json:"id"`
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

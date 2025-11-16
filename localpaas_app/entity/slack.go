package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type Slack struct {
	Webhook string `json:"webhook"`

	// NOTE: for storing current containing setting only
	Setting *Setting `json:"-"`
}

func (s *Setting) ParseSlack() (*Slack, error) {
	res := &Slack{Setting: s}
	if s != nil && s.Data != "" {
		err := s.parseData(res)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return res, nil
	}
	return res, nil
}

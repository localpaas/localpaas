package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type Slack struct {
	Webhook string `json:"webhook"`
}

func (s *Setting) ParseSlack() (*Slack, error) {
	if s != nil && s.Data != "" {
		res := &Slack{}
		err := s.parseData(res)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return res, nil
	}
	return nil, nil
}

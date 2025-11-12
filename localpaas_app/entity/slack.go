package entity

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
)

type Slack struct {
	Webhook string `json:"webhook"`
}

func (s *Setting) ParseSlack() (*Slack, error) {
	res := &Slack{}
	if s != nil && s.Data != "" {
		err := s.parseData(res)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		return res, nil
	}
	return res, nil
}

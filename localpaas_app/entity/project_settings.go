package entity

import (
	"encoding/json"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/pkg/reflectutil"
)

type ProjectSettings struct {
	Test string `json:"test"`
}

func (p *Project) GetSettings() (*ProjectSettings, error) {
	if len(p.AllSettings) == 0 {
		return nil, nil
	}
	data := p.AllSettings[0].Data
	if len(data) == 0 {
		return nil, nil
	}
	res := &ProjectSettings{}
	err := json.Unmarshal(reflectutil.UnsafeStrToBytes(data), res)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return res, nil
}

type ProjectEnvVars struct {
	Data [][]string `json:"data"`
}

func (p *Project) GetEnvVars() (*ProjectEnvVars, error) {
	if len(p.AllSettings) == 0 {
		return nil, nil
	}
	data := p.AllSettings[0].Data
	if len(data) == 0 {
		return nil, nil
	}
	res := &ProjectEnvVars{}
	err := json.Unmarshal(reflectutil.UnsafeStrToBytes(data), res)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return res, nil
}

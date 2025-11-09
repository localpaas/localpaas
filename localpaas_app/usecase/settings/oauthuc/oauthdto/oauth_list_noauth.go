package oauthdto

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type ListOAuthNoAuthReq struct {
	Name   []string             `json:"-"`
	Status []base.SettingStatus `json:"-"`
	Search string               `json:"-"`

	Paging basedto.Paging `json:"-"`
}

type ListOAuthNoAuthResp struct {
	Meta *basedto.Meta `json:"meta"`
	Data []*OAuthResp  `json:"data"`
}

package sessiondto

import "github.com/localpaas/localpaas/localpaas_app/basedto"

type RefreshSessionResp struct {
	Meta *basedto.BaseMeta       `json:"meta"`
	Data *RefreshSessionDataResp `json:"data"`
}

type RefreshSessionDataResp struct {
	*BaseCreateSessionResp
}

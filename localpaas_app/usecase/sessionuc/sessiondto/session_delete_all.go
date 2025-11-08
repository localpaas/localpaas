package sessiondto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type DeleteAllSessionsReq struct {
	User *basedto.User `json:"-"`
}

func NewDeleteAllSessionsReq() *DeleteAllSessionsReq {
	return &DeleteAllSessionsReq{}
}

func (req *DeleteAllSessionsReq) Validate() apperrors.ValidationErrors {
	return nil
}

type DeleteAllSessionsResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

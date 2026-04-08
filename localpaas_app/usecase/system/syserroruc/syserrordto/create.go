package syserrordto

import (
	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateSysErrorReq struct {
	ErrorInfo *apperrors.ErrorInfo `json:"-"`
}

func NewCreateSysErrorReq() *CreateSysErrorReq {
	return &CreateSysErrorReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *CreateSysErrorReq) Validate() apperrors.ValidationErrors {
	return nil
}

type CreateSysErrorResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

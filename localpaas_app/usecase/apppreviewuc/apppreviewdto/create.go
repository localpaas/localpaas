package apppreviewdto

import (
	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

const (
	pullRequestMaxLen = 100
)

type CreatePreviewReq struct {
	ProjectID string `json:"-"`
	AppID     string `json:"-"`

	PullRequest string `json:"pullRequest"`
}

func NewCreatePreviewReq() *CreatePreviewReq {
	return &CreatePreviewReq{}
}

func (req *CreatePreviewReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, basedto.ValidateID(&req.ProjectID, true, "projectId")...)
	validators = append(validators, basedto.ValidateID(&req.AppID, true, "appId")...)
	validators = append(validators, basedto.ValidateStr(&req.PullRequest, true, 1, pullRequestMaxLen,
		"pullRequest")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreatePreviewResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

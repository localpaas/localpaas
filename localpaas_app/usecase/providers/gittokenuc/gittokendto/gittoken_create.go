package gittokendto

import (
	"strings"
	"time"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type CreateGitTokenReq struct {
	*GitTokenBaseReq
}

type GitTokenBaseReq struct {
	Kind     base.GitSource `json:"kind"`
	Name     string         `json:"name"`
	User     string         `json:"user"`
	Token    string         `json:"token"`
	ExpireAt time.Time      `json:"expireAt"`
}

func (req *GitTokenBaseReq) modifyRequest() error {
	req.Name = strings.TrimSpace(req.Name)
	req.User = strings.TrimSpace(req.User)
	req.Token = strings.TrimSpace(req.Token)
	return nil
}

func (req *GitTokenBaseReq) validate(_ string) []vld.Validator {
	// TODO: add validation
	return nil
}

func NewCreateGitTokenReq() *CreateGitTokenReq {
	return &CreateGitTokenReq{}
}

func (req *CreateGitTokenReq) ModifyRequest() error {
	return req.modifyRequest()
}

// Validate implements interface basedto.ReqValidator
func (req *CreateGitTokenReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type CreateGitTokenResp struct {
	Meta *basedto.BaseMeta     `json:"meta"`
	Data *basedto.ObjectIDResp `json:"data"`
}

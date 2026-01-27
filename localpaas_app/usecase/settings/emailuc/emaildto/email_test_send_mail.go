package emaildto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestSendEmailReq struct {
	*EmailBaseReq
	TestRecipient string `json:"testRecipient"`
	TestSubject   string `json:"testSubject"`
	TestContent   string `json:"testContent"`
}

func NewTestSendEmailReq() *TestSendEmailReq {
	return &TestSendEmailReq{}
}

func (req *TestSendEmailReq) ModifyRequest() error {
	// NOTE: make sure req.Name is not empty to not fail the validation
	req.Name = gofn.Coalesce(req.Name, "x")
	req.TestSubject = gofn.Coalesce(req.TestSubject, "test subject")
	req.TestContent = gofn.Coalesce(req.TestContent, "test content")
	return nil
}

// Validate implements interface basedto.ReqValidator
func (req *TestSendEmailReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestSendEmailResp struct {
	Meta *basedto.Meta `json:"meta"`
}

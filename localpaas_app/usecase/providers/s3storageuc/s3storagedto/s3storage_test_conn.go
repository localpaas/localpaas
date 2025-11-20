package s3storagedto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

type TestS3StorageConnReq struct {
	*S3StorageBaseReq
}

func NewTestS3StorageConnReq() *TestS3StorageConnReq {
	return &TestS3StorageConnReq{}
}

// Validate implements interface basedto.ReqValidator
func (req *TestS3StorageConnReq) Validate() apperrors.ValidationErrors {
	// NOTE: make sure req.Name is not empty to not fail the validation
	req.Name = gofn.Coalesce(req.Name, "x")

	var validators []vld.Validator
	validators = append(validators, req.validate("")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type TestS3StorageConnResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

package filedto

import (
	"io"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings"
)

type DownloadFileReq struct {
	settings.GetSettingReq
	Token                   string            `json:"-" mapstructure:"token"`
	ViewInline              bool              `json:"-" mapstructure:"viewInline"`
	UsePresignURLOnFileSize int64             `json:"-" mapstructure:"-"`
	PresignExpiration       timeutil.Duration `json:"-" mapstructure:"-"`
}

func NewDownloadFileReq() *DownloadFileReq {
	return &DownloadFileReq{}
}

func (req *DownloadFileReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, req.GetSettingReq.Validate()...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type DownloadFileResp struct {
	Meta *basedto.Meta         `json:"meta"`
	Data *DownloadFileDataResp `json:"data"`
}

type DownloadFileDataResp struct {
	RedirectURL   string            `json:"redirectURL"`
	ContentType   string            `json:"contentType"`
	ContentLength int64             `json:"contentLength"`
	ExtraHeaders  map[string]string `json:"headers"`
	Content       io.ReadCloser     `json:"content"`
}

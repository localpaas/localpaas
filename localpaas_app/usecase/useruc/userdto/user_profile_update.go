package userdto

import (
	"encoding/base64"
	"path/filepath"
	"strings"

	vld "github.com/tiendc/go-validator"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/pkg/strutil"
)

type UpdateProfileReq struct {
	Username string        `json:"username"`
	Email    string        `json:"email"`
	FullName string        `json:"fullName"`
	Photo    *UserPhotoReq `json:"photo"`
}

func NewUpdateProfileReq() *UpdateProfileReq {
	return &UpdateProfileReq{}
}

func (req *UpdateProfileReq) ModifyRequest() error {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strutil.NormalizeEmail(req.Email)
	req.FullName = strings.TrimSpace(req.FullName)
	// Parse photo
	if req.Photo != nil && req.Photo.FileName != "" && req.Photo.DataBase64 != "" {
		req.Photo.DataBytes, _ = base64.StdEncoding.DecodeString(req.Photo.DataBase64)
		req.Photo.FileExt = strings.ToLower(filepath.Ext(req.Photo.FileName))
	}
	return nil
}

func (req *UpdateProfileReq) Validate() apperrors.ValidationErrors {
	var validators []vld.Validator
	validators = append(validators, validateUsername(&req.Username, false, "username")...)
	validators = append(validators, basedto.ValidateEmail(&req.Email, false, "email")...)
	validators = append(validators, basedto.ValidateStr(&req.FullName, false,
		minNameLen, maxNameLen, "fullName")...)
	validators = append(validators, validateUserPhoto(req.Photo, "photo")...)
	return apperrors.NewValidationErrors(vld.Validate(validators...))
}

type UpdateProfileResp struct {
	Meta *basedto.BaseMeta `json:"meta"`
}

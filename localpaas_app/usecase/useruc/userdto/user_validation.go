package userdto

import (
	"github.com/asaskevich/govalidator"
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
)

var (
	acceptedPhotoFileExts = []string{".png", ".jpg", ".jpeg", ".webp"}
)

const (
	minNameLen = 1
	maxNameLen = 100

	maxUserPhotoSize = 300 * 1024 // 300KB

	minNotesLen = 1
	maxNotesLen = 10000
)

func validateUsername(username *string, required bool, field string) (res []vld.Validator) {
	res = append(res, basedto.ValidateStr(username, required,
		minNameLen, maxNameLen, "username")...)

	// NOTE: username must not be a valid email address as it takes the address
	if username != nil {
		res = append(res, vld.Must(!govalidator.IsEmail(*username)).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_PARAM_INVALID"),
		))
	}
	return res
}

func validateUserPhoto(photo *UserPhotoReq, field string) []vld.Validator {
	if photo == nil || photo.FileName == "" || photo.DataBase64 == "" {
		return nil
	}
	return []vld.Validator{
		vld.Must(gofn.Contain(acceptedPhotoFileExts, photo.FileExt)).OnError(
			vld.SetField(field+".fileName", nil),
			vld.SetCustomKey("ERR_VLD_USER_PHOTO_FILE_EXT_UNSUPPORTED"),
		),
		vld.Must(len(photo.DataBytes) > 0).OnError(
			vld.SetField(field+".dataBase64", nil),
			vld.SetCustomKey("ERR_VLD_USER_PHOTO_FILE_INVALID"),
		),
		vld.When(len(photo.DataBytes) > 0).Then(
			vld.Must(len(photo.DataBytes) <= maxUserPhotoSize).OnError(
				vld.SetField(field+".dataBase64", nil),
				vld.SetCustomKey("ERR_VLD_USER_PHOTO_FILE_TOO_BIG"),
			),
		),
	}
}

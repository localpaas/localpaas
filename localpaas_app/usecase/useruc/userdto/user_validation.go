package userdto

import (
	vld "github.com/tiendc/go-validator"
	"github.com/tiendc/gofn"
)

var (
	acceptedPhotoFileExts = []string{".png", ".jpg", ".jpeg", ".webp"}
)

const (
	minNameLen = 1
	maxNameLen = 100

	maxUserPhotoSize = 500 * 1024 * 1024 // 500KB
)

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

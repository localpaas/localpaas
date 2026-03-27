package fileservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
	"github.com/localpaas/localpaas/localpaas_app/service/settingservice"
)

type FileService interface {
	GetDownloadURL(ctx context.Context, db database.IDB, auth *basedto.Auth,
		req *GetDownloadURLReq) (*GetDownloadURLResp, error)
}

func NewFileService(
	settingRepo repository.SettingRepo,
	settingService settingservice.SettingService,
) FileService {
	return &fileService{
		settingRepo:    settingRepo,
		settingService: settingService,
	}
}

type fileService struct {
	settingRepo    repository.SettingRepo
	settingService settingservice.SettingService
}

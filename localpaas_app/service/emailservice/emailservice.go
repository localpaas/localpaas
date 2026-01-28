package emailservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/repository"
)

type EmailService interface {
	GetDefaultSystemEmail(ctx context.Context, db database.IDB) (*entity.Setting, error)

	SendMailPasswordReset(ctx context.Context, db database.IDB, data *EmailDataPasswordReset) error
}

func NewEmailService(
	settingRepo repository.SettingRepo,
) EmailService {
	return &emailService{
		settingRepo: settingRepo,
	}
}

type emailService struct {
	settingRepo repository.SettingRepo
}

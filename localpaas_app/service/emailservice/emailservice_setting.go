package emailservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
)

func (s *emailService) GetDefaultSystemEmail(
	ctx context.Context,
	db database.IDB,
) (*entity.Setting, error) {
	settings, _, err := s.settingRepo.List(ctx, db, nil,
		bunex.SelectWhere("setting.type = ?", base.SettingTypeEmail),
		bunex.SelectWhere("setting.object_id IS NULL"),
		bunex.SelectOrder("setting.is_default DESC"),
		bunex.SelectLimit(1),
	)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	if len(settings) == 0 {
		return nil, apperrors.NewNotFound("EmailSetting").WithMsgLog("default system email not found")
	}
	return settings[0], nil
}

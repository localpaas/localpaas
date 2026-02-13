package settingservice

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
)

type CreateEvent struct {
	Setting *entity.Setting
}

func (s *settingService) OnCreate(
	ctx context.Context,
	db database.IDB,
	event *CreateEvent,
) (err error) {
	return nil
}

type UpdateEvent struct {
	Setting *entity.Setting

	OldStatus base.SettingStatus
	NewStatus base.SettingStatus
	OldKind   string
	NewKind   string
}

func (s *settingService) OnUpdate(
	ctx context.Context,
	db database.IDB,
	event *UpdateEvent,
) (err error) {
	// Remove healthcheck cache if the update may relate
	if event.Setting.IsTypeIn(base.SettingTypeHealthcheck, base.SettingTypeIMService, base.SettingTypeEmail) {
		err = s.healthcheckSettingsRepo.Del(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

type DeleteEvent struct {
	Setting *entity.Setting
}

func (s *settingService) OnDelete(
	ctx context.Context,
	db database.IDB,
	event *DeleteEvent,
) (err error) {
	// Remove healthcheck cache if the update may relate
	if event.Setting.IsTypeIn(base.SettingTypeHealthcheck, base.SettingTypeIMService, base.SettingTypeEmail) {
		err = s.healthcheckSettingsRepo.Del(ctx)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

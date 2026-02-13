package settingservice

import (
	"context"

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
) error {
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
) error {
	return nil
}

type DeleteEvent struct {
	Setting *entity.Setting
}

func (s *settingService) OnDelete(
	ctx context.Context,
	db database.IDB,
	event *DeleteEvent,
) error {
	return nil
}

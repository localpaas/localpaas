package settingserviceimpl

import (
	"context"
	"time"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/ulid"
)

func (s *service) initDefaultNotificationSettings(
	ctx context.Context,
	db database.IDB,
	timeNow time.Time,
) (err error) {
	notifSetting := &entity.Setting{
		ID:              gofn.Must(ulid.NewStringULID()),
		Scope:           base.ObjectScopeGlobal,
		Type:            base.SettingTypeNotification,
		Status:          base.SettingStatusActive,
		Name:            "default",
		AvailInProjects: false,
		Default:         true,
		Version:         entity.CurrentNotificationVersion,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
	}
	notifSetting.MustSetData(entity.NewNotificationDefaultForScope(base.NewObjectScopeGlobal()))

	err = s.settingRepo.Insert(ctx, db, notifSetting)
	if err != nil {
		return apperrors.New(err)
	}
	return nil
}

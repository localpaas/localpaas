package settingservice

import "github.com/localpaas/localpaas/localpaas_app/entity"

type PersistingSettingData struct {
	UpsertingSettings []*entity.Setting
	UpsertingAccesses []*entity.ACLPermission
	DeletingAccesses  []*entity.ACLPermission
}

type CreateEvent struct {
	Setting *entity.Setting
}

type DeleteEvent struct {
	Setting *entity.Setting
}

type UpdateEvent struct {
	Setting    *entity.Setting
	OldSetting *entity.Setting
}

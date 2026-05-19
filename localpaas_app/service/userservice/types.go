package userservice

import (
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	SkipCheckingCurrentPassword = ""
)

type PersistingUserData struct {
	UpsertingUsers      []*entity.User
	UpsertingSettings   []*entity.Setting
	UpsertingBinObjects []*entity.BinObject
	UpsertingAccesses   []*entity.ACLPermission
	DeletingAccesses    []*base.PermissionResource
}

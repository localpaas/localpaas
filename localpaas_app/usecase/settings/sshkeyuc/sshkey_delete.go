package sshkeyuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/settings/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) DeleteSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.DeleteSSHKeyReq,
) (*sshkeydto.DeleteSSHKeyResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		sshKeyData := &deleteSSHKeyData{}
		err := uc.loadSSHKeyDataForDelete(ctx, db, req, sshKeyData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSSHKeyData{}
		uc.prepareDeletingSSHKey(sshKeyData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.DeleteSSHKeyResp{}, nil
}

type deleteSSHKeyData struct {
	Setting *entity.Setting
}

func (uc *SSHKeyUC) loadSSHKeyDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *sshkeydto.DeleteSSHKeyReq,
	data *deleteSSHKeyData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSSHKey, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *SSHKeyUC) prepareDeletingSSHKey(
	data *deleteSSHKeyData,
	persistingData *persistingSSHKeyData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}

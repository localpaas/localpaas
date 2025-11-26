package sshkeyuc

import (
	"context"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/sshkeyuc/sshkeydto"
)

func (uc *SSHKeyUC) UpdateSSHKey(
	ctx context.Context,
	auth *basedto.Auth,
	req *sshkeydto.UpdateSSHKeyReq,
) (*sshkeydto.UpdateSSHKeyResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		sshKeyData := &updateSSHKeyData{}
		err := uc.loadSSHKeyDataForUpdate(ctx, db, req, sshKeyData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSSHKeyData{}
		err = uc.prepareUpdatingSSHKey(req.SSHKeyPartialReq, sshKeyData, persistingData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &sshkeydto.UpdateSSHKeyResp{}, nil
}

type updateSSHKeyData struct {
	Setting *entity.Setting
}

func (uc *SSHKeyUC) loadSSHKeyDataForUpdate(
	ctx context.Context,
	db database.IDB,
	req *sshkeydto.UpdateSSHKeyReq,
	data *updateSSHKeyData,
) error {
	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSSHKey, req.ID, false,
		bunex.SelectFor("UPDATE OF setting"),
		bunex.SelectRelation("ObjectAccesses",
			bunex.SelectWhere("acl_permission.subject_type IN (?)", bunex.In([]base.SubjectType{
				base.SubjectTypeProject, base.SubjectTypeApp,
			})),
		),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	if req.UpdateVer != setting.UpdateVer {
		return apperrors.Wrap(apperrors.ErrUpdateVerMismatched)
	}
	data.Setting = setting

	// If name changes, validate the new one
	if req.Name != nil && !strings.EqualFold(setting.Name, *req.Name) {
		conflictSetting, _ := uc.settingRepo.GetByName(ctx, db, base.SettingTypeSSHKey, *req.Name, false)
		if conflictSetting != nil {
			return apperrors.NewAlreadyExist("SSHKey").
				WithMsgLog("ssh key '%s' already exists", conflictSetting.Name)
		}
	}

	return nil
}

func (uc *SSHKeyUC) prepareUpdatingSSHKey(
	req *sshkeydto.SSHKeyPartialReq,
	data *updateSSHKeyData,
	persistingData *persistingSSHKeyData,
) error {
	timeNow := timeutil.NowUTC()
	setting := data.Setting
	if req.Name != nil {
		setting.Name = *req.Name
	}

	if req.PrivateKey != nil {
		sshKey, err := setting.AsSSHKey()
		if err != nil {
			return apperrors.Wrap(err)
		}
		if sshKey == nil {
			sshKey = &entity.SSHKey{}
		}
		if req.PrivateKey != nil {
			sshKey.PrivateKey = *req.PrivateKey
		}
		if req.Passphrase != nil {
			sshKey.Passphrase = *req.Passphrase
		}

		setting.MustSetData(sshKey.MustEncrypt())
	}

	setting.UpdatedAt = timeNow
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)

	// Project accesses change
	if req.ProjectAccesses != nil {
		// Remove all current items
		persistingData.DeletingAccesses = append(persistingData.DeletingAccesses, setting.ObjectAccesses...)
		uc.preparePersistingSSHKeyProjects(setting, req.ProjectAccesses, timeNow, persistingData)
	}
	return nil
}
